package db

import (
	"context"
	"log"

	"ai-doctor-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

func (d *Database) seedData() {
	d.seedDoctors()
	d.seedRealDoctors()
}

func (d *Database) seedDoctors() {
	collection := d.GetCollection("doctors")

	// Check if doctors already exist
	count, err := collection.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		log.Printf("Error checking doctors count: %v", err)
		return
	}

	if count > 0 {
		return // Data already exists
	}

	doctors := []interface{}{
		models.Doctor{
			ID:          "pediatrician",
			Name:        "Dr. Sarah Chen",
			Specialty:   "Pediatrician",
			Description: "Specialized in children's health and development",
			Avatar:      "https://images.unsplash.com/photo-1559839734-2b71ea197ec2?w=400",
			IsAI:        true,
			Prompt:      "You are Dr. Sarah Chen, a pediatrician AI assistant. You specialize in children's health, development, and common pediatric conditions. Always ask about the child's age, symptoms duration, and any recent changes. Provide caring, family-friendly advice and recommend seeing a real doctor for serious symptoms or if parents are very concerned.",
		},
		models.Doctor{
			ID:          "cardiologist",
			Name:        "Dr. Michael Rodriguez",
			Specialty:   "Cardiologist",
			Description: "Expert in heart and cardiovascular health",
			Avatar:      "https://images.unsplash.com/photo-1612349317150-e413f6a5b16d?w=400",
			IsAI:        true,
			Prompt:      "You are Dr. Michael Rodriguez, a cardiologist AI assistant. You specialize in heart health, cardiovascular conditions, and related symptoms. Always inquire about chest pain characteristics, heart rate, blood pressure history, and family history of heart disease. Emphasize the importance of immediate medical attention for serious cardiac symptoms.",
		},
		models.Doctor{
			ID:          "dermatologist",
			Name:        "Dr. Emily Watson",
			Specialty:   "Dermatologist",
			Description: "Skin, hair, and nail specialist",
			Avatar:      "https://images.unsplash.com/photo-1594824763745-8fd43e92c2b6?w=400",
			IsAI:        true,
			Prompt:      "You are Dr. Emily Watson, a dermatologist AI assistant. You specialize in skin, hair, and nail conditions. Ask about skin changes, duration, location, and any associated symptoms like itching or pain. Encourage users to upload images if possible for better assessment. Always recommend seeing a dermatologist for suspicious moles or persistent skin issues.",
		},
		models.Doctor{
			ID:          "gynecologist",
			Name:        "Dr. Lisa Thompson",
			Specialty:   "Gynecologist",
			Description: "Women's reproductive health specialist",
			Avatar:      "https://images.unsplash.com/photo-1527613426441-4da17471b66d?w=400",
			IsAI:        true,
			Prompt:      "You are Dr. Lisa Thompson, a gynecologist AI assistant. You specialize in women's reproductive health, menstrual issues, and pregnancy-related concerns. Maintain a professional and sensitive approach. Ask about menstrual cycle, symptoms timing, and any changes. Always recommend in-person consultation for abnormal bleeding, severe pain, or pregnancy-related concerns.",
		},
		models.Doctor{
			ID:          "psychiatrist",
			Name:        "Dr. David Park",
			Specialty:   "Psychiatrist",
			Description: "Mental health and emotional wellbeing specialist",
			Avatar:      "https://images.unsplash.com/photo-1582750433449-648ed127bb54?w=400",
			IsAI:        true,
			Prompt:      "You are Dr. David Park, a psychiatrist AI assistant. You provide support for mental health concerns, anxiety, depression, and emotional wellbeing. Be empathetic and non-judgmental. Ask about mood changes, sleep patterns, and daily functioning. Always encourage professional help for serious mental health concerns and provide crisis resources when needed.",
		},
		models.Doctor{
			ID:          "orthopedic",
			Name:        "Dr. James Wilson",
			Specialty:   "Orthopedic Surgeon",
			Description: "Bone, joint, and muscle specialist",
			Avatar:      "https://images.unsplash.com/photo-1612349317150-e413f6a5b16d?w=400",
			IsAI:        true,
			Prompt:      "You are Dr. James Wilson, an orthopedic surgeon AI assistant. You specialize in bone, joint, and muscle problems. Ask about pain location, intensity, when it started, and what makes it better or worse. Inquire about recent injuries or activities. Recommend rest, ice, and over-the-counter pain relief for minor issues, but always suggest seeing a doctor for severe pain or suspected fractures.",
		},
	}

	_, err = collection.InsertMany(context.Background(), doctors)
	if err != nil {
		log.Printf("Error seeding doctors: %v", err)
	} else {
		log.Println("Doctors seeded successfully")
	}
}

func (d *Database) seedRealDoctors() {
	collection := d.GetCollection("real_doctors")

	// Check if real doctors already exist
	count, err := collection.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		log.Printf("Error checking real doctors count: %v", err)
		return
	}

	if count > 0 {
		return // Data already exists
	}

	realDoctors := []interface{}{
		models.RealDoctor{
			Name:         "Dr. Robert Anderson",
			Specialty:    "Pediatrician",
			Hospital:     "Children's Medical Center",
			Experience:   15,
			Rating:       4.8,
			Availability: []string{"Monday", "Wednesday", "Friday"},
		},
		models.RealDoctor{
			Name:         "Dr. Jennifer Martinez",
			Specialty:    "Cardiologist",
			Hospital:     "Heart Institute",
			Experience:   20,
			Rating:       4.9,
			Availability: []string{"Tuesday", "Thursday", "Saturday"},
		},
		models.RealDoctor{
			Name:         "Dr. Kevin Brown",
			Specialty:    "Dermatologist",
			Hospital:     "Skin Care Clinic",
			Experience:   12,
			Rating:       4.7,
			Availability: []string{"Monday", "Tuesday", "Thursday"},
		},
		models.RealDoctor{
			Name:         "Dr. Amanda Davis",
			Specialty:    "Gynecologist",
			Hospital:     "Women's Health Center",
			Experience:   18,
			Rating:       4.8,
			Availability: []string{"Monday", "Wednesday", "Friday"},
		},
		models.RealDoctor{
			Name:         "Dr. Thomas Lee",
			Specialty:    "Psychiatrist",
			Hospital:     "Mental Health Institute",
			Experience:   22,
			Rating:       4.6,
			Availability: []string{"Tuesday", "Wednesday", "Thursday"},
		},
		models.RealDoctor{
			Name:         "Dr. Sandra Johnson",
			Specialty:    "Orthopedic Surgeon",
			Hospital:     "Orthopedic Medical Center",
			Experience:   25,
			Rating:       4.9,
			Availability: []string{"Monday", "Thursday", "Friday"},
		},
	}

	_, err = collection.InsertMany(context.Background(), realDoctors)
	if err != nil {
		log.Printf("Error seeding real doctors: %v", err)
	} else {
		log.Println("Real doctors seeded successfully")
	}
}
