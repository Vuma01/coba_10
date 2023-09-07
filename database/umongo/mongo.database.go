package umongo

import (
	"coba_01/pkg/helper"
	"coba_01/src/config"
	"coba_01/src/models"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

// MongoDBUtil merupakan struktur yang menyediakan detail
// untuk koneksi dan operasi terhadap MongoDB.
type MongoDBUtil struct {
	server         string          // Alamat server MongoDB
	dbname         string          // Nama database yang digunakan
	collectionName string          // Nama koleksi dalam database
	ctx            context.Context // Konteks untuk operasi database yang bersifat async
}

// NewMongoDBUtil adalah fungsi pembuat untuk MongoDBUtil.
// Fungsi ini menginisialisasi MongoDBUtil dengan informasi yang diperlukan
// untuk berkomunikasi dengan database, termasuk mengambil informasi
// server dan nama database dari variabel lingkungan dan
// menyetel nama koleksi dari parameter yang diberikan.
func NewMongoDBUtil(collectionName string) *MongoDBUtil {
	return &MongoDBUtil{
		server:         os.Getenv(config.ENV_MONGO_SRV), // Ambil alamat server MongoDB dari variabel lingkungan
		dbname:         os.Getenv(config.ENV_DB_NAME),   // Ambil nama database dari variabel lingkungan
		collectionName: collectionName,                  // Set nama koleksi dari parameter
		ctx:            context.Background(),            // Set konteks background untuk operasi async
	}
}

// ConnDB adalah metode pada MongoDBUtil yang bertugas untuk membuat
// koneksi ke server MongoDB dan mengembalikan klien MongoDB.
// Jika ada kesalahan saat mencoba menyambungkan ke server MongoDB,
// ErrorHelper akan dipanggil untuk menangani kesalahan tersebut.
func (db MongoDBUtil) ConnDB() (client *mongo.Client, err error) {
	clientOptions := options.Client()                  // Buat opsi klien baru untuk konfigurasi koneksi
	clientOptions.ApplyURI(db.server)                  // Tetapkan URI server menggunakan alamat server yang diberikan dalam struktur
	client, err = mongo.Connect(db.ctx, clientOptions) // Buat koneksi ke server MongoDB menggunakan opsi yang diberikan
	helper.ErrorHelperPanic(err)                       // Jika ada kesalahan, panggil ErrorHelper untuk menangani
	return
}

// IsEmailExist memeriksa apakah sebuah email sudah ada di dalam database.
func (db MongoDBUtil) IsEmailExist(email string) (bool, error) {
	// Membuat koneksi ke database.
	client, err := db.ConnDB()
	if err != nil {
		return false, err
	}
	// Pastikan untuk menutup koneksi setelah selesai.
	defer client.Disconnect(db.ctx)

	// Mengambil koleksi dari database berdasarkan nama koleksi yang telah ditentukan.
	coll := client.Database(db.dbname).Collection(db.collectionName)

	// Menghitung jumlah dokumen yang memiliki email yang sesuai.
	count, err := coll.CountDocuments(db.ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}

	// Mengembalikan true jika ditemukan lebih dari 0 dokumen.
	return count > 0, nil
}

// Insert memasukkan data ke dalam koleksi database.
func (db MongoDBUtil) Insert(data interface{}) error {
	// Membuat koneksi ke database.
	client, err := db.ConnDB()
	if err != nil {
		return err
	}
	// Pastikan untuk menutup koneksi setelah selesai.
	defer client.Disconnect(db.ctx)

	// Mengambil koleksi dan memasukkan data.
	coll := client.Database(db.dbname).Collection(db.collectionName)
	_, err = coll.InsertOne(context.TODO(), data)
	return err
}

// GetAllWithPaging mengambil semua data user dengan paginasi.
func (db MongoDBUtil) GetAllWithPaging(skip int, limit int) ([]models.User, error) {
	// Membuat koneksi ke database.
	client, err := db.ConnDB()
	if err != nil {
		return nil, err
	}
	// Pastikan untuk menutup koneksi setelah selesai.
	defer client.Disconnect(db.ctx)

	// Mengambil koleksi dan melakukan pencarian dengan paginasi.
	coll := client.Database(db.dbname).Collection(db.collectionName)
	cursor, err := coll.Find(db.ctx, bson.M{}, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}

	// Memuat hasil ke dalam slice users.
	var users []models.User
	if err = cursor.All(db.ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (db MongoDBUtil) IsEmailUsedByAnotherUser(email string, usedId primitive.ObjectID) (bool, error) {
	client, err := db.ConnDB()
	if err != nil {
		return false, err
	}
	defer client.Disconnect(db.ctx)

	coll := client.Database(db.dbname).Collection(db.collectionName)
	filter := bson.M{"email": email, "_id": bson.M{"$ne": usedId}}
	count, err := coll.CountDocuments(db.ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (db MongoDBUtil) Update(filter, update bson.M) error {
	client, err := db.ConnDB()
	if err != nil {
		return err
	}
	defer client.Disconnect(db.ctx)

	coll := client.Database(db.dbname).Collection(db.collectionName)
	_, err = coll.UpdateOne(db.ctx, filter, update)
	return err
}

func (db MongoDBUtil) Delete(id string) (int64, error) {
	client, err := db.ConnDB()
	if err != nil {
		return 0, err
	}
	defer client.Disconnect(db.ctx)

	coll := client.Database(db.dbname).Collection(db.collectionName)

	objID, _ := primitive.ObjectIDFromHex(id)
	result, err := coll.DeleteOne(db.ctx, bson.M{"_id": objID})
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

// ================>[ Pencarian ID - Email ]<================\\

func (db MongoDBUtil) GetUserByID(id string) (*models.User, error) {
	client, err := db.ConnDB()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(db.ctx)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	col := client.Database(db.dbname).Collection(db.collectionName)
	var user models.User
	if err = col.FindOne(db.ctx, bson.M{"_id": objID}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(fmt.Sprintf("ID = %s %s ", id, " Tidak Ditemukan"))
		}
		return nil, err
	}
	return &user, nil
}
func (db MongoDBUtil) GetUserByEmail(email string) ([]models.User, error) {
	client, err := db.ConnDB()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(db.ctx)

	coll := client.Database(db.dbname).Collection(db.collectionName)

	var user []models.User
	filter := bson.M{"email": bson.M{"$regex": email, "$options": "i"}}
	cursor, err := coll.Find(db.ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(db.ctx)
	if err = cursor.All(db.ctx, &user); err != nil {
		return nil, err
	}
	return user, nil
}
func (db MongoDBUtil) GetEmailLogin(email string) (*models.User, error) {
	client, err := db.ConnDB()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(db.ctx)

	coll := client.Database(db.dbname).Collection(db.collectionName)
	var user models.User

	filter := bson.M{"email": email}

	err = coll.FindOne(db.ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
