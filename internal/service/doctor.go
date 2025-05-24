package service

import (
	"context"

	"github.com/subhammahanty235/medai/internal/db"
	"github.com/subhammahanty235/medai/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

type DoctorService struct {
	db *db.Database
}

func NewDoctorService(database *db.Database) *DoctorService {
	return &DoctorService{
		db: database,
	}
}

func (s *DoctorService) GetAllDoctors() ([]models.Doctor, error) {
	collection := s.db.GetCollection("doctors")

	cursor, err := collection.Find(context.Background(), bson.M{"is_ai": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var doctors []models.Doctor
	if err = cursor.All(context.Background(), &doctors); err != nil {
		return nil, err
	}

	return doctors, nil
}

func (s *DoctorService) GetDoctorByID(doctorID string) (*models.Doctor, error) {
	collection := s.db.GetCollection("doctors")

	var doctor models.Doctor
	err := collection.FindOne(context.Background(), bson.M{"_id": doctorID}).Decode(&doctor)
	if err != nil {
		return nil, err
	}

	return &doctor, nil
}

func (s *DoctorService) GetRealDoctorsBySpecialty(specialty string) ([]models.RealDoctor, error) {
	collection := s.db.GetCollection("real_doctors")

	cursor, err := collection.Find(context.Background(), bson.M{"specialty": specialty})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var doctors []models.RealDoctor
	if err = cursor.All(context.Background(), &doctors); err != nil {
		return nil, err
	}

	return doctors, nil
}

func (s *DoctorService) GetAllRealDoctors() ([]models.RealDoctor, error) {
	collection := s.db.GetCollection("real_doctors")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var doctors []models.RealDoctor
	if err = cursor.All(context.Background(), &doctors); err != nil {
		return nil, err
	}

	return doctors, nil
}
