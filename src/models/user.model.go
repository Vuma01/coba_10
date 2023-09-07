package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// User adalah model untuk pengguna.
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username" json:"username" validate:"required,min=4"`
	Email     string             `bson:"email" json:"email" validate:"required,email"`
	Password  string             `bson:"password" json:"password" validate:"required,min=8"`
	Address   Address            `bson:"address" json:"address"`
	Role      string             `bson:"role" json:"role" validate:"required,oneof=user admin"` // Misalnya: user, admin
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	OTPInfo   OTPInfo            `bson:"otp_info" json:"otp_info"`
}

// Address adalah model untuk alamat pengguna.
type Address struct {
	Kota string `bson:"kota" json:"kota"`
	Pos  string `bson:"pos" json:"pos"`
}

// OTPInfo berisi informasi tentang OTP untuk reset password atau verifikasi.
type OTPInfo struct {
	Code         string    `bson:"code" json:"code"`                   // Kode OTP
	GeneratedAt  time.Time `bson:"generated_at" json:"generated_at"`   // Waktu ketika OTP di-generate
	ExpiryTime   time.Time `bson:"expiry_time" json:"expiry_time"`     // Waktu kedaluwarsa OTP
	IsUsed       bool      `bson:"is_used" json:"is_used"`             // Menandai apakah OTP sudah digunakan atau belum
	LastTriedAt  time.Time `bson:"last_tried_at" json:"last_tried_at"` // Waktu terakhir pengguna mencoba memasukkan OTP
	AttemptCount int       `bson:"attempt_count" json:"attempt_count"` // Jumlah percobaan pengguna memasukkan OTP
}
