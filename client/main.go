package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

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

type IP_Asset struct {
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
	Hyperlink          string `json:"hyperlink"`
	Status             string `json:"status"`
	Certificate        string `json:"certificate"`
}

type Publication struct {
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
	NumberOfCitations    int    `json:"number_of_citation"`
	GoogleScholarDetails string `json:"google_scholar_details"`
	SDGNo                string `json:"sdg_no"`
	FundingType          string `json:"funding_type"`
	NatureOfFunding      string `json:"nature_of_funding"`
	Publisher            string `json:"publisher"`
	Abstract             string `json:"abstract"`
}

type User struct {
	UserID      int    `json:"user_id"`
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

type Log struct {
	LogID       string `json:"log_id"`
	Activity    string `json:"activity"`
	Description string `json:"description"`
}

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	client := pb.NewRMSServiceClient(conn)

	r := gin.Default()

	//table_authors
	r.GET("/table_authors", func(ctx *gin.Context) {
		res, err := client.GetAuthors(ctx, &pb.ReadAuthorsRequest{})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"table_authors": res.Authors,
		})
	})
	r.GET("/table_authors/:author_id", func(ctx *gin.Context) {
		id := ctx.Param("author_id")
		res, err := client.GetAuthor(ctx, &pb.ReadAuthorRequest{AuthorId: id})
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"table_authors": res.Author,
		})
	})
	r.POST("/table_authors", func(ctx *gin.Context) {
		var author Author

		err := ctx.ShouldBind(&author)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		data := &pb.Author{
			AuthorId:     author.ID,
			AuthorName:   author.AuthorName,
			Gender:       author.AuthorGender,
			TypeOfAuthor: author.TypeofAuthor,
			Affiliation:  author.Affiliation,
			Email:        author.AuthorEmail,
		}
		res, err := client.CreateAuthor(ctx, &pb.CreateAuthorRequest{
			Author: data,
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"": res.Author,
		})
	})
	r.PUT("/table_authors/:author_id", func(ctx *gin.Context) {
		var author Author
		err := ctx.ShouldBind(&author)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		res, err := client.UpdateAuthor(ctx, &pb.UpdateAuthorRequest{
			Author: &pb.Author{
				AuthorId:     author.ID,
				AuthorName:   author.AuthorName,
				Gender:       author.AuthorGender,
				TypeOfAuthor: author.TypeofAuthor,
				Affiliation:  author.Affiliation,
				Email:        author.AuthorEmail,
			},
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"table_authors": res.Author,
		})
		return

	})
	r.DELETE("/table_authors/:author_id", func(ctx *gin.Context) {
		id := ctx.Param("author_id")
		res, err := client.DeleteAuthor(ctx, &pb.DeleteAuthorRequest{AuthorId: id})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if res.Success == true {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Author deleted successfully",
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "error deleting author",
			})
			return
		}

	})

	//table_ipassets
	r.GET("/table_ipassets", func(ctx *gin.Context) {
		res, err := client.GetIP_Assets(ctx, &pb.ReadIP_AssetsRequest{})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"table_ipassets": res.IpAssets,
		})
	})
	r.GET("/table_ipassets/:registration_number", func(ctx *gin.Context) {
		id := ctx.Param("registration_number")
		res, err := client.GetIP_Asset(ctx, &pb.ReadIP_AssetRequest{RegistrationNumber: id})
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"table_ipassets": res.IpAsset,
		})
	})
	r.POST("/table_ipassets", func(ctx *gin.Context) {
		var ipAsset IP_Asset

		err := ctx.ShouldBind(&ipAsset)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		data := &pb.IP_Asset{
			RegistrationNumber: ipAsset.RegistartionNumber,
			TitleOfWork:        ipAsset.TitleOfWork,
			TypeOfDocument:     ipAsset.TypeOfDocument,
			ClassOfWork:        ipAsset.ClassOfWork,
			DateOfCreation:     ipAsset.DateOfCreation,
			DateRegistered:     ipAsset.DateRegistered,
			Campus:             ipAsset.Campus,
			College:            ipAsset.College,
			Program:            ipAsset.Program,
			Authors:            ipAsset.Authors,
			Hyperlink:          ipAsset.Hyperlink,
			Status:             ipAsset.Status,
			Certificate:        ipAsset.Certificate,
		}
		res, err := client.CreateIP_Asset(ctx, &pb.CreateIP_AssetRequest{
			IpAsset: data,
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"table_ipassets": res.IpAsset,
		})
	})
	r.PUT("/table_ipassets/:registration_number", func(ctx *gin.Context) {
		var ipAsset IP_Asset
		err := ctx.ShouldBind(&ipAsset)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		res, err := client.UpdateIP_Asset(ctx, &pb.UpdateIP_AssetRequest{
			IpAsset: &pb.IP_Asset{
				RegistrationNumber: ipAsset.RegistartionNumber,
				TitleOfWork:        ipAsset.TitleOfWork,
				TypeOfDocument:     ipAsset.TypeOfDocument,
				ClassOfWork:        ipAsset.ClassOfWork,
				DateOfCreation:     ipAsset.DateOfCreation,
				DateRegistered:     ipAsset.DateRegistered,
				Campus:             ipAsset.Campus,
				College:            ipAsset.College,
				Program:            ipAsset.Program,
				Authors:            ipAsset.Authors,
				Hyperlink:          ipAsset.Hyperlink,
				Status:             ipAsset.Status,
				Certificate:        ipAsset.Certificate,
			},
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"table_ipassets": res.IpAsset,
		})
		return

	})
	r.DELETE("/table_ipassets/:registration_number", func(ctx *gin.Context) {
		id := ctx.Param("registration_number")
		res, err := client.DeleteIP_Asset(ctx, &pb.DeleteIP_AssetRequest{RegistrationNumber: id})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if res.Success == true {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "IP asset deleted successfully",
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "error deleting ip asset",
			})
			return
		}

	})

	//table_publications
	r.GET("/table_publications", func(ctx *gin.Context) {
		res, err := client.GetPublications(ctx, &pb.ReadPublicationsRequest{})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"table_publications": res.Publications,
		})
	})
	r.GET("/table_publications/:publication_id", func(ctx *gin.Context) {
		id := ctx.Param("publication_id")
		res, err := client.GetPublication(ctx, &pb.ReadPublicationRequest{PublicationId: id})
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"table_publications": res.Publication,
		})
	})
	r.POST("/table_publications", func(ctx *gin.Context) {
		var publication Publication

		err := ctx.ShouldBind(&publication)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		data := &pb.Publication{
			PublicationId:        publication.PublicationID,
			DatePublished:        publication.DatePublished,
			Quartile:             publication.Quartile,
			Authors:              publication.Authors,
			Department:           publication.Department,
			College:              publication.College,
			Campus:               publication.Campus,
			TitleOfPaper:         publication.TitleOfPaper,
			TypeOfPublication:    publication.TypeOfPublication,
			FundingSource:        publication.FundingSource,
			NumberOfCitation:     int32(publication.NumberOfCitations),
			GoogleScholarDetails: publication.GoogleScholarDetails,
			SdgNo:                publication.SDGNo,
			FundingType:          publication.FundingType,
			NatureOfFunding:      publication.NatureOfFunding,
			Publisher:            publication.Publisher,
			Abstract:             publication.Abstract,
		}
		res, err := client.CreatePublication(ctx, &pb.CreatePublicationRequest{
			Publication: data,
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"table_publications": res.Publication,
		})
	})
	r.PUT("/table_publications/:publication_id", func(ctx *gin.Context) {
		var publication Publication
		err := ctx.ShouldBind(&publication)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		numberOfCitations, _ := strconv.Atoi(ctx.Param("number_of_citation"))
		res, err := client.UpdatePublication(ctx, &pb.UpdatePublicationRequest{
			Publication: &pb.Publication{
				PublicationId:        publication.PublicationID,
				DatePublished:        publication.DatePublished,
				Quartile:             publication.Quartile,
				Authors:              publication.Authors,
				Department:           publication.Department,
				College:              publication.College,
				Campus:               publication.Campus,
				TitleOfPaper:         publication.TitleOfPaper,
				TypeOfPublication:    publication.TypeOfPublication,
				FundingSource:        publication.FundingSource,
				NumberOfCitation:     int32(numberOfCitations),
				GoogleScholarDetails: publication.GoogleScholarDetails,
				SdgNo:                publication.SDGNo,
				FundingType:          publication.FundingType,
				NatureOfFunding:      publication.NatureOfFunding,
				Publisher:            publication.Publisher,
				Abstract:             publication.Abstract,
			},
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"table_publications": res.Publication,
		})
		return

	})
	r.DELETE("/table_publications/:publication_id", func(ctx *gin.Context) {
		id := ctx.Param("publication_id")
		res, err := client.DeletePublication(ctx, &pb.DeletePublicationRequest{PublicationId: id})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if res.Success == true {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Publication deleted successfully",
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "error deleting publication",
			})
			return
		}

	})

	//table_user
	r.GET("/table_user", func(ctx *gin.Context) {
		res, err := client.GetUsers(ctx, &pb.ReadUsersRequest{})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"table_user": res.Users,
		})
	})
	r.GET("/table_user/:user_id", func(ctx *gin.Context) {
		id := ctx.Param("user_id")
		userID, _ := strconv.Atoi(id)
		res, err := client.GetUser(ctx, &pb.ReadUserRequest{UserId: int32(userID)})
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"table_user": res.User,
		})
	})
	r.POST("/table_user", func(ctx *gin.Context) {
		var user User

		err := ctx.ShouldBind(&user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		data := &pb.User{
			UserId:      int32(user.UserID),
			SrCode:      user.SRCode,
			Email:       user.Email,
			Password:    user.Password,
			AccountType: user.AccountType,
			UserContact: user.UserContact,
			UserImg:     user.UserImg,
			UserFname:   user.UserFname,
			UserLname:   user.UserLname,
			UserMname:   user.UserMname,
		}
		res, err := client.CreateUser(ctx, &pb.CreateUserRequest{
			User: data,
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"table_user": res.User,
		})
	})
	r.PUT("/table_user/:user_id", func(ctx *gin.Context) {
		var user User
		err := ctx.ShouldBind(&user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		userID, _ := strconv.Atoi(ctx.Param("user_id"))
		res, err := client.UpdateUser(ctx, &pb.UpdateUserRequest{
			User: &pb.User{
				UserId:      int32(userID),
				SrCode:      user.SRCode,
				Email:       user.Email,
				Password:    user.Password,
				AccountType: user.AccountType,
				UserContact: user.UserContact,
				UserImg:     user.UserImg,
				UserFname:   user.UserFname,
				UserLname:   user.UserLname,
				UserMname:   user.UserMname,
			},
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"table_user": res.User,
		})
		return

	})
	r.DELETE("/table_user/:user_id", func(ctx *gin.Context) {
		id := ctx.Param("user_id")
		userID, _ := strconv.Atoi(id)
		res, err := client.DeleteUser(ctx, &pb.DeleteUserRequest{UserId: int32(userID)})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if res.Success == true {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "User deleted successfully",
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "error deleting user",
			})
			return
		}

	})

	//table_log
	r.GET("/table_log", func(ctx *gin.Context) {
		res, err := client.GetLogs(ctx, &pb.ReadLogsRequest{})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"table_log": res.Logs,
		})
	})
	r.GET("/table_log/:log_id", func(ctx *gin.Context) {
		id := ctx.Param("log_id")
		res, err := client.GetLog(ctx, &pb.ReadLogRequest{LogId: id})
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"table_log": res.Log,
		})
	})
	r.POST("/table_log", func(ctx *gin.Context) {
		var log Log

		err := ctx.ShouldBind(&log)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		data := &pb.Log{
			LogId:       log.LogID,
			Activity:    log.Activity,
			Description: log.Description,
		}
		res, err := client.CreateLog(ctx, &pb.CreateLogRequest{
			Log: data,
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"table_log": res.Log,
		})
	})
	r.PUT("/table_log/:log_id", func(ctx *gin.Context) {
		var log Log
		err := ctx.ShouldBind(&log)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		res, err := client.UpdateLog(ctx, &pb.UpdateLogRequest{
			Log: &pb.Log{
				LogId:       log.LogID,
				Activity:    log.Activity,
				Description: log.Description,
			},
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"table_log": res.Log,
		})
		return

	})
	r.DELETE("/table_log/:log_id", func(ctx *gin.Context) {
		id := ctx.Param("log_id")
		res, err := client.DeleteLog(ctx, &pb.DeleteLogRequest{LogId: id})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if res.Success == true {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Log deleted successfully",
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "error deleting log",
			})
			return
		}

	})

	r.Run(":5000")

}
