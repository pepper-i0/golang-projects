package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/pepper/ass-db/models"
	"github.com/pepper/ass-db/storage"
	"gorm.io/gorm"
)

type Admin struct {
	Firstname string `json: "firstname"`
	Lastname  string `json: "lastname"`
	Email     string `json: "email"`
	Password  string `json: "password"`
	Role      string `json: "role"`
}

type User struct {
	Firstname string `json: "firstname"`
	Lastname  string `json: "lastname"`
	Email     string `json: "email"`
	Phone     int    `json: "phone"`
}
type Repository struct {
	DB *gorm.DB
}

func (r *Repository) GetAdmins(context *fiber.Ctx) error {
	adminModels := &[]models.Admins{}

	err := r.DB.Find(adminModels).Error
	if err != nil {
		context.Status(http.StatusBadGateway).JSON(
			&fiber.Map{"message": "could not find admin"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Admins fetched successfully",
			"data": adminModels,
		})
	return nil
}

func (r *Repository) GetAdminByID(context *fiber.Ctx) error {
	id := context.Params("id")
	adminModel := &models.Admins{}

	if id == "" {
		context.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Map{"message": "id cannot be empty"})
		return nil
	}

	fmt.Println("The id is", id)

	err := r.DB.Where("id = ?", id).First(adminModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			fiber.Map{"message": "could not get admin"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"messsage": "user id fetched successfully",
			"data": adminModel,
		})
	return nil
}

func (r *Repository) CreateAdmin(context *fiber.Ctx) error {
	admin := Admin{}

	err := context.BodyParser(&admin)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "could not create admin"})
		return err
	}

	err = r.DB.Create(&admin).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create admin"})
		return err

	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "admin added successfully"})

	return nil
}

func (r *Repository) DeleteAdmin(context *fiber.Ctx) error {
	adminModel := models.Admins{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "id cannot be empty"})
		return nil
	}
	err := r.DB.Delete(adminModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not delete admin"})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "admin deleted successfully"})
	return nil
}

func (r *Repository) CreatePost(context *fiber.Ctx) error {
	post := &models.Posts{}
	err := context.BodyParser(&post)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Invalid request body"})
		return err
	}

	pendingStatus := string(models.Pending)
	post.ReviewStatus = &pendingStatus

	err = r.DB.Create(&post).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create post"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "post created successfully", "data": post})
	return nil
}

func (r *Repository) EditPost(context *fiber.Ctx) error {
	id := context.Params("id")
	post := &models.Posts{}
	err := r.DB.First(post, id).Error
	if err != nil {
		context.Status(http.StatusNotFound).JSON(
			&fiber.Map{"message": "Post not found"})
		return nil
	}

	err = context.BodyParser(post)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Invalid request body"})
		return err
	}

	err = r.DB.Save(post).Error
	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "could not update post"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Post updateed Successfully", "data": post})
	return nil
}

func (r *Repository) ReviewPost(context *fiber.Ctx) error {
	id := context.Params("id")
	post := &models.Posts{}
	err := r.DB.First(post, id).Error
	if err != nil {
		context.Status(http.StatusNotFound).JSON(
			fiber.Map{"message": "Post not found"})
		return err
	}

	var input struct {
		Reviewstatus string `json:"review_status"`
		AdminID      uint   `json:"admin_id"`
	}

	err = context.BodyParser(&input)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "invalid request body"})
		return err
	}

	if input.Reviewstatus != "approved" && input.Reviewstatus != "declined" {
		context.Status(http.StatusBadGateway).JSON(
			&fiber.Map{"message": "invalid review status"})
		return err
	}
	post.ReviewStatus = &input.Reviewstatus
	post.AdminID = &input.AdminID
	post.ReviewDate = time.Now()

	err = r.DB.Save(post).Error
	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "could not update post"})
		return err
	}
	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Post reviewed successfully", "data": post})
	return nil
}

func (r *Repository) GetPendingPosts(context *fiber.Ctx) error {
	posts := &[]models.Posts{}
	err := r.DB.Where("ReviewStatus=?", "pending").Find(&posts).Error
	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "couln not fetch pending posts"})
		return err
	}
	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Pending posts fetched successfully"})
	return nil
}

func (r *Repository) CreateUser(context *fiber.Ctx) error {
	user := User{}

	err := context.BodyParser(&user)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"messge": "request failed"})
		return err
	}

	err = r.DB.Create(&user).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not create user"})
		return err
	}
	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "User has been added"})
	return nil
}

func (r *Repository) DeleteUser(context *fiber.Ctx) error {
	userModel := models.Users{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "id cannot be empty"})
		return nil
	}
	err := r.DB.Delete(userModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "counld not delete user"})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "user deleted successfully"})
	return nil
}

func (r *Repository) GetPosts(context *fiber.Ctx) error {
	postModels := &[]models.Posts{}
	err := r.DB.Find(&postModels).Error
	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "could not fetch  posts"})
		return err
	}
	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Posts fetched successfully",
			"data": postModels,
		})
	return nil
}

func (r *Repository) GetUsers(context *fiber.Ctx) error {
	userModels := &[]models.Users{}

	err := r.DB.Find(userModels).Error
	if err != nil {
		context.Status(http.StatusBadGateway).JSON(
			&fiber.Map{"message": "could not find user"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "users fetched successfully",
			"data": userModels,
		})
	return nil
}

func (r *Repository) GetUserByID(context *fiber.Ctx) error {
	id := context.Params("id")
	userModel := &models.Users{}
	if id == "" {
		context.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Map{"message": "id cannot be empty"})
		return nil
	}

	fmt.Println("The id is", id)

	err := r.DB.Where("id = ?", id).First(userModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get user"})
		return err
	}
	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "user id fetched successfully",
			"data": userModel,
		})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/Create_admin", r.CreateAdmin)
	api.Delete("/Delete_admin/:id", r.DeleteAdmin)
	api.Get("Get_admin/:id", r.GetAdminByID)
	api.Get("/Admins", r.GetAdmins)
	api.Post("/Create_user", r.CreateUser)
	api.Delete("Delete_user/:id", r.DeleteUser)
	api.Get("/Get_user/:id", r.GetUserByID)
	api.Get("/Users", r.GetUsers)
	api.Post("/Create_Posts", r.CreatePost)
	api.Patch("/posts/:id", r.EditPost)
	api.Patch("/posts/:id/review", r.ReviewPost)
	api.Get("/posts/pending", r.GetPendingPosts)
	api.Post("/posts", r.GetPosts)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODEL"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Could not load Database")
	}
	err = models.MigrateUsers(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	err = models.MigratePosts(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	err = models.MigrateAdmins(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")

}
