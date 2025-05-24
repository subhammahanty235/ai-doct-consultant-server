package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type Doctor struct {
	ID          string `bson:"_id" json:"id"`
	Name        string `bson:"name" json:"name"`
	Specialty   string `bson:"specialty" json:"specialty"`
	Description string `bson:"description" json:"description"`
	Avatar      string `bson:"avatar" json:"avatar"`
	IsAI        bool   `bson:"is_ai" json:"is_ai"`
	Prompt      string `bson:"prompt" json:"-"`
}

type RealDoctor struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name         string             `bson:"name" json:"name"`
	Specialty    string             `bson:"specialty" json:"specialty"`
	Hospital     string             `bson:"hospital" json:"hospital"`
	Experience   int                `bson:"experience" json:"experience"`
	Rating       float64            `bson:"rating" json:"rating"`
	Availability []string           `bson:"availability" json:"availability"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

type ChatSession struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	DoctorID  string             `bson:"doctor_id" json:"doctor_id"`
	Messages  []Message          `bson:"messages" json:"messages"`
	Status    string             `bson:"status" json:"status"` // active, completed, doctor_recommended
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Content   string             `bson:"content" json:"content"`
	Sender    string             `bson:"sender" json:"sender"` // user, ai, system
	ImageURL  string             `bson:"image_url,omitempty" json:"image_url,omitempty"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}

type Appointment struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID           primitive.ObjectID `bson:"user_id" json:"user_id"`
	RealDoctorID     primitive.ObjectID `bson:"real_doctor_id" json:"real_doctor_id"`
	ChatSessionID    primitive.ObjectID `bson:"chat_session_id" json:"chat_session_id"`
	AppointmentDate  time.Time          `bson:"appointment_date" json:"appointment_date"`
	Status           string             `bson:"status" json:"status"` // pending, confirmed, completed, cancelled
	Symptoms         string             `bson:"symptoms" json:"symptoms"`
	AIRecommendation string             `bson:"ai_recommendation" json:"ai_recommendation"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}

// Request/Response DTOs
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type ChatMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

type AppointmentRequest struct {
	RealDoctorID     string    `json:"real_doctor_id" binding:"required"`
	AppointmentDate  time.Time `json:"appointment_date" binding:"required"`
	Symptoms         string    `json:"symptoms"`
	AIRecommendation string    `json:"ai_recommendation"`
}
