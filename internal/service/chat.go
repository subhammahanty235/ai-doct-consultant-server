package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ai-doctor-backend/internal/db"
	"ai-doctor-backend/internal/models"
	"ai-doctor-backend/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatService struct {
	db            *db.Database
	geminiClient  *utils.GeminiClient
	doctorService *DoctorService
}

func NewChatService(database *db.Database, geminiClient *utils.GeminiClient, doctorService *DoctorService) *ChatService {
	return &ChatService{
		db:            database,
		geminiClient:  geminiClient,
		doctorService: doctorService,
	}
}

func (s *ChatService) StartChatSession(userID primitive.ObjectID, doctorID string) (*models.ChatSession, error) {
	collection := s.db.GetCollection("chat_sessions")

	// Check if there's an active session for this user and doctor
	var existingSession models.ChatSession
	err := collection.FindOne(context.Background(), bson.M{
		"user_id":   userID,
		"doctor_id": doctorID,
		"status":    "active",
	}).Decode(&existingSession)

	if err == nil {
		return &existingSession, nil // Return existing active session
	}

	// Create new session
	session := models.ChatSession{
		UserID:    userID,
		DoctorID:  doctorID,
		Messages:  []models.Message{},
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := collection.InsertOne(context.Background(), session)
	if err != nil {
		return nil, err
	}

	session.ID = result.InsertedID.(primitive.ObjectID)

	// Add welcome message
	doctor, err := s.doctorService.GetDoctorByID(doctorID)
	if err != nil {
		return nil, err
	}

	welcomeMessage := models.Message{
		ID:        primitive.NewObjectID(),
		Content:   fmt.Sprintf("Hello! I'm %s, your AI %s. How can I help you today? Please tell me about your symptoms or concerns.", doctor.Name, doctor.Specialty),
		Sender:    "ai",
		Timestamp: time.Now(),
	}

	session.Messages = append(session.Messages, welcomeMessage)

	// Update session with welcome message
	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{"_id": session.ID},
		bson.M{"$set": bson.M{"messages": session.Messages, "updated_at": time.Now()}},
	)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *ChatService) SendMessage(sessionID primitive.ObjectID, userID primitive.ObjectID, content string, imageURL string) (*models.Message, error) {
	collection := s.db.GetCollection("chat_sessions")

	// Get session
	var session models.ChatSession
	err := collection.FindOne(context.Background(), bson.M{"_id": sessionID, "user_id": userID}).Decode(&session)
	if err != nil {
		return nil, err
	}

	// Add user message
	userMessage := models.Message{
		ID:        primitive.NewObjectID(),
		Content:   content,
		Sender:    "user",
		ImageURL:  imageURL,
		Timestamp: time.Now(),
	}

	session.Messages = append(session.Messages, userMessage)

	// Get doctor info for AI response
	doctor, err := s.doctorService.GetDoctorByID(session.DoctorID)
	if err != nil {
		return nil, err
	}

	// Build conversation context
	conversationContext := s.buildConversationContext(session.Messages)
	fullPrompt := fmt.Sprintf("%s\n\nConversation so far:\n%s\n\nLatest user message: %s", doctor.Prompt, conversationContext, content)

	// Generate AI response
	var aiResponse string
	if imageURL != "" {
		aiResponse, err = s.geminiClient.GenerateResponseWithImage(doctor.Prompt, fullPrompt, imageURL)
	} else {
		aiResponse, err = s.geminiClient.GenerateResponse(doctor.Prompt, fullPrompt)
	}

	if err != nil {
		return nil, err
	}

	// Check if AI recommends seeing a real doctor
	shouldRecommendDoctor := s.shouldRecommendRealDoctor(aiResponse, content)
	if shouldRecommendDoctor {
		session.Status = "doctor_recommended"
		aiResponse += "\n\n🏥 Based on your symptoms, I recommend scheduling an appointment with a real doctor for proper examination and treatment. Would you like me to help you find available doctors in my specialty?"
	}

	// Add AI response
	aiMessage := models.Message{
		ID:        primitive.NewObjectID(),
		Content:   aiResponse,
		Sender:    "ai",
		Timestamp: time.Now(),
	}

	session.Messages = append(session.Messages, aiMessage)

	// Update session
	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{"_id": sessionID},
		bson.M{"$set": bson.M{
			"messages":   session.Messages,
			"status":     session.Status,
			"updated_at": time.Now(),
		}},
	)
	if err != nil {
		return nil, err
	}

	return &aiMessage, nil
}

func (s *ChatService) GetChatHistory(userID primitive.ObjectID) ([]models.ChatSession, error) {
	collection := s.db.GetCollection("chat_sessions")

	cursor, err := collection.Find(
		context.Background(),
		bson.M{"user_id": userID},
		// Sort by most recent first
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var sessions []models.ChatSession
	if err = cursor.All(context.Background(), &sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (s *ChatService) GetChatSession(sessionID primitive.ObjectID, userID primitive.ObjectID) (*models.ChatSession, error) {
	collection := s.db.GetCollection("chat_sessions")

	var session models.ChatSession
	err := collection.FindOne(context.Background(), bson.M{"_id": sessionID, "user_id": userID}).Decode(&session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *ChatService) buildConversationContext(messages []models.Message) string {
	var context strings.Builder
	for _, msg := range messages {
		if msg.Sender == "user" {
			context.WriteString(fmt.Sprintf("Patient: %s\n", msg.Content))
		} else if msg.Sender == "ai" {
			context.WriteString(fmt.Sprintf("AI Doctor: %s\n", msg.Content))
		}
	}
	return context.String()
}

func (s *ChatService) shouldRecommendRealDoctor(aiResponse, userMessage string) bool {
	// Simple keyword-based detection for recommending real doctor
	concerningKeywords := []string{
		"severe", "emergency", "urgent", "chest pain", "difficulty breathing",
		"blood", "seizure", "unconscious", "broken bone", "fracture",
		"suicidal", "heart attack", "stroke", "high fever", "persistent pain",
	}

	responseKeywords := []string{
		"recommend seeing", "consult a doctor", "medical attention",
		"see a specialist", "hospital", "emergency room",
	}

	userLower := strings.ToLower(userMessage)
	responseLower := strings.ToLower(aiResponse)

	for _, keyword := range concerningKeywords {
		if strings.Contains(userLower, keyword) {
			return true
		}
	}

	for _, keyword := range responseKeywords {
		if strings.Contains(responseLower, keyword) {
			return true
		}
	}

	return false
}
