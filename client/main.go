package main

import (
	"flag"
	"log"
	"net/http"

	pb "example.com/go-grpc-crud-api/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

type Author struct {
	ID           string `json:"author_id"`
	AuthorName   string `json:"author_name"`
	AuthorGender string `json:"gender"`
	TypeofAuthor string `json:"type_of_author"`
	Affiliation  string `json:"affiliation"`
	AuthorEmail  string `json:"email"`
}

type IP_Assets struct {
	RegistartionNumber string `json:"registration_number"`
	TitleOfWork        string `json:"title_of_work"`
	TypeOfDocument     string `json:"type_of_document"`
	ClassOfWork        string `json:"class_of_work"`
	DateOfCreation     string `json:"date_of_creation"`
	DateRegistered     string `json:"date_registered"`
	Campus             string `json:"campus"`
	College            string `json:"college"`
	Program            string `json:"program"`
	Authors            string `json:"authors"`
}

type Publications struct {
	PublicationID        string `json:"publication_id"`
	DatePublished        string `json:"date_published"`
	Quartile             string `json:"quartile"`
	Authors              string `json:"authors"`
	Department           string `json:"department"`
	College              string `json:"college"`
	Campus               string `json:"campus"`
	TitleOfPaper         string `json:"title_of_paper"`
	TypeOfPublication    string `json:"type_of_publication"`
	FundingSource        string `json:"funding_source"`
	NumberOfCitations    string `json:"number_of_citations"`
	GoogleScholarDetails string `json:"google_scholar_details"`
	SDGNo                string `json:"sdg_no"`
	FundingType          string `json:"funding_type"`
	NatureOfFundings     string `json:"nature_of_fundings"`
	Publisher            string `json:"publisher"`
	Abstract             string `json:"abstract"`
}

type User struct {
	UserID      int64  `json:"user_id"`
	SRCode      string `json:"sr_code"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	AccountType string `json:"account_type"`
	UserContact string `json:"user_contact"`
	UserImg     string `json:"user_img"`
	UserFname   string `json:"user_fname"`
	UserLname   string `json:"user_lname"`
	UserMname   string `json:"user_mname"`
}

type Movie struct {
	ID    string `json:"id"`
	Title string `json:"Title"`
	Genre string `json:"genre"`
}

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	client := pb.NewMovieServiceClient(conn)

	r := gin.Default()
	r.GET("/movies", func(ctx *gin.Context) {
		res, err := client.GetMovies(ctx, &pb.ReadMoviesRequest{})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"movies": res.Movies,
		})
	})
	r.GET("/movies/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		res, err := client.GetMovie(ctx, &pb.ReadMovieRequest{Id: id})
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"movie": res.Movie,
		})
	})
	r.POST("/movies", func(ctx *gin.Context) {
		var movie Movie

		err := ctx.ShouldBind(&movie)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		data := &pb.Movie{
			Title: movie.Title,
			Genre: movie.Genre,
		}
		res, err := client.CreateMovie(ctx, &pb.CreateMovieRequest{
			Movie: data,
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"movie": res.Movie,
		})
	})
	r.PUT("/movies/:id", func(ctx *gin.Context) {
		var movie Movie
		err := ctx.ShouldBind(&movie)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		res, err := client.UpdateMovie(ctx, &pb.UpdateMovieRequest{
			Movie: &pb.Movie{
				Id:    movie.ID,
				Title: movie.Title,
				Genre: movie.Genre,
			},
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"movie": res.Movie,
		})
		return

	})
	r.DELETE("/movies/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		res, err := client.DeleteMovie(ctx, &pb.DeleteMovieRequest{Id: id})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if res.Success == true {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Movie deleted successfully",
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "error deleting movie",
			})
			return
		}

	})

	r.Run(":5000")

}
