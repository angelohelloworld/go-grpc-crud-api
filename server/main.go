package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "example.com/go-grpc-crud-api/proto"
	"google.golang.org/grpc"

	"github.com/google/uuid"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	DatabaseConnection()
}

var DB *gorm.DB
var err error

type Author struct {
	ID           string `gorm:"primarykey"`
	AuthorName   string
	AuthorGender string
	TypeofAuthor string
	Affiliation  string
	AuthorEmail  string
	CreatedAt    time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime:false"`
}

type IP_Asset struct {
	RegistartionNumber string `gorm:"primarykey"`
	TitleOfWork        string
	TypeOfDocument     string
	ClassOfWork        string
	DateOfCreation     string
	DateRegistered     string
	Campus             string
	College            string
	Program            string
	Authors            string
	CreatedAt          time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime:false"`
}

type Publication struct {
	PublicationID        string `gorm:"primarykey"`
	DatePublished        string
	Quartile             string
	Authors              string
	Department           string
	College              string
	Campus               string
	TitleOfPaper         string
	TypeOfPublication    string
	FundingSource        string
	NumberOfCitations    string
	GoogleScholarDetails string
	SDGNo                string
	FundingType          string
	NatureOfFundings     string
	Publisher            string
	Abstract             string
	CreatedAt            time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime:false"`
}

type User struct {
	UserID      int32 `gorm:"primarykey"`
	SRCode      string
	Email       string
	Password    string
	AccountType string
	UserContact string
	UserImg     string
	UserFname   string
	UserLname   string
	UserMname   string
	CreatedAt   time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime:false"`
}

func DatabaseConnection() {
	host := "localhost"
	port := "5432"
	dbName := "RMS_db"
	dbUser := "postgres"
	password := "password"
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		dbUser,
		dbName,
		password,
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database...", err)
	}

	DB.AutoMigrate(&Author{})
	DB.AutoMigrate(&IP_Asset{})
	DB.AutoMigrate(&Publication{})
	DB.AutoMigrate(&User{})

	fmt.Println("Database connection successful...")
}

var (
	port = flag.Int("port", 50051, "gRPC server port")
)

type server struct {
	pb.UnimplementedRMSServiceServer
}

// Author
func (*server) CreateAuthor(ctx context.Context, req *pb.CreateAuthorRequest) (*pb.CreateAuthorResponse, error) {
	fmt.Println("Create Author")
	author := req.GetAuthor()
	author.AuthorId = uuid.New().String()

	data := Author{

		ID:           author.GetAuthorId(),
		AuthorName:   author.GetAuthorName(),
		AuthorGender: author.GetGender(),
		TypeofAuthor: author.GetTypeOfAuthor(),
		Affiliation:  author.GetAffiliation(),
		AuthorEmail:  author.GetEmail(),
	}

	res := DB.Table("table_authors").Create(&data)
	if res.RowsAffected == 0 {
		return nil, errors.New("author creation unsuccessful")
	}
	return &pb.CreateAuthorResponse{
		Author: &pb.Author{
			AuthorId:     author.GetAuthorId(),
			AuthorName:   author.GetAuthorName(),
			Gender:       author.GetGender(),
			TypeOfAuthor: author.GetTypeOfAuthor(),
			Affiliation:  author.GetAffiliation(),
			Email:        author.GetEmail(),
		},
	}, nil
}

func (*server) GetAuthor(ctx context.Context, req *pb.ReadAuthorRequest) (*pb.ReadAuthorResponse, error) {
	fmt.Println("Read Author", req.GetAuthorId())
	var author Author
	res := DB.Table("table_authors").Find(&author, "Authorid = ?", req.GetAuthorId())
	if res.RowsAffected == 0 {
		return nil, errors.New("Author not found")
	}
	return &pb.ReadAuthorResponse{
		Author: &pb.Author{
			AuthorId:     author.ID,
			AuthorName:   author.AuthorName,
			Gender:       author.AuthorGender,
			TypeOfAuthor: author.TypeofAuthor,
			Affiliation:  author.Affiliation,
			Email:        author.AuthorEmail,
		},
	}, nil
}

func (*server) GetAuthors(ctx context.Context, req *pb.ReadAuthorsRequest) (*pb.ReadAuthorsResponse, error) {
	fmt.Println("Read Authors")
	authors := []*pb.Author{}
	res := DB.Table("table_authors").Find(&authors)
	if res.RowsAffected == 0 {
		return nil, errors.New("author not found")
	}
	return &pb.ReadAuthorsResponse{
		Authors: authors,
	}, nil
}

func (*server) UpdateAuthor(ctx context.Context, req *pb.UpdateAuthorRequest) (*pb.UpdateAuthorResponse, error) {
	fmt.Println("Update Author")
	var author Author
	reqAuthor := req.GetAuthor()

	res := DB.Table("table_authors").Model(&author).Where("Authorid=?", reqAuthor.AuthorId).Updates(
		Author{
			AuthorName:   reqAuthor.AuthorName,
			AuthorGender: reqAuthor.Gender,
			TypeofAuthor: reqAuthor.TypeOfAuthor,
			Affiliation:  reqAuthor.Affiliation,
			AuthorEmail:  reqAuthor.Email,
		})

	if res.RowsAffected == 0 {
		return nil, errors.New("Author not found")
	}

	return &pb.UpdateAuthorResponse{
		Author: &pb.Author{
			AuthorId:     author.ID,
			AuthorName:   author.AuthorName,
			Gender:       author.AuthorGender,
			TypeOfAuthor: author.TypeofAuthor,
			Affiliation:  author.Affiliation,
			Email:        author.AuthorEmail,
		},
	}, nil
}

func (*server) DeleteAuthor(ctx context.Context, req *pb.DeleteAuthorRequest) (*pb.DeleteAuthorResponse, error) {
	fmt.Println("Delete Author")
	var author Author
	res := DB.Table("table_authors").Where("Authorid = ?", req.GetAuthorId()).Delete(&author)
	if res.RowsAffected == 0 {
		return nil, errors.New("author not found")
	}

	return &pb.DeleteAuthorResponse{
		Success: true,
	}, nil
}

// IP_Asset
func (*server) CreateIPAsset(ctx context.Context, req *pb.CreateIP_AssetRequest) (*pb.CreateIP_AssetResponse, error) {
	fmt.Println("Create IP Asset")
	ipAsset := req.GetIpAsset()
	ipAsset.RegistrationNumber = uuid.New().String()

	data := IP_Asset{
		RegistartionNumber: ipAsset.GetRegistrationNumber(),
		TitleOfWork:        ipAsset.GetTitleOfWork(),
		TypeOfDocument:     ipAsset.GetTypeOfDocument(),
		ClassOfWork:        ipAsset.GetClassOfWork(),
		DateOfCreation:     ipAsset.GetDateOfCreation(),
		DateRegistered:     ipAsset.GetDateRegistered(),
		Campus:             ipAsset.GetCampus(),
		College:            ipAsset.GetCollege(),
		Program:            ipAsset.GetProgram(),
		Authors:            ipAsset.GetAuthors(),
	}

	res := DB.Table("table_ipassets").Create(&data)
	if res.RowsAffected == 0 {
		return nil, errors.New("IP asset creation unsuccessful")
	}

	return &pb.CreateIP_AssetResponse{
		IpAsset: &pb.IP_Asset{
			RegistrationNumber: ipAsset.GetRegistrationNumber(),
			TitleOfWork:        ipAsset.GetTitleOfWork(),
			TypeOfDocument:     ipAsset.GetTypeOfDocument(),
			ClassOfWork:        ipAsset.GetClassOfWork(),
			DateOfCreation:     ipAsset.GetDateOfCreation(),
			DateRegistered:     ipAsset.GetDateRegistered(),
			Campus:             ipAsset.GetCampus(),
			College:            ipAsset.GetCollege(),
			Program:            ipAsset.GetProgram(),
			Authors:            ipAsset.GetAuthors(),
		},
	}, nil
}

func (*server) GetIPAsset(ctx context.Context, req *pb.ReadIP_AssetRequest) (*pb.ReadIP_AssetResponse, error) {
	fmt.Println("Read IP Asset", req.GetRegistrationNumber())
	var ipAsset IP_Asset
	res := DB.Table("table_ipassets").Find(&ipAsset, "RegistartionNumber = ?", req.GetRegistrationNumber())
	if res.RowsAffected == 0 {
		return nil, errors.New("IP asset not found")
	}
	return &pb.ReadIP_AssetResponse{
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
		},
	}, nil
}

func (*server) GetIPAssets(ctx context.Context, req *pb.ReadIP_AssetsRequest) (*pb.ReadIP_AssetsResponse, error) {
	fmt.Println("Read IP Assets")
	ipAssets := []*pb.IP_Asset{}
	res := DB.Table("table_ipassets").Find(&ipAssets)
	if res.RowsAffected == 0 {
		return nil, errors.New("IP assets not found")
	}
	return &pb.ReadIP_AssetsResponse{
		IpAssets: ipAssets,
	}, nil
}

func (*server) UpdateIPAsset(ctx context.Context, req *pb.UpdateIP_AssetRequest) (*pb.UpdateIP_AssetResponse, error) {
	fmt.Println("Update IP Asset")
	var ipAsset IP_Asset
	reqIPAsset := req.GetIpAsset()

	res := DB.Table("table_ipassets").Model(&ipAsset).Where("RegistartionNumber = ?", reqIPAsset.RegistrationNumber).Updates(
		IP_Asset{
			TitleOfWork:    reqIPAsset.TitleOfWork,
			TypeOfDocument: reqIPAsset.TypeOfDocument,
			ClassOfWork:    reqIPAsset.ClassOfWork,
			DateOfCreation: reqIPAsset.DateOfCreation,
			DateRegistered: reqIPAsset.DateRegistered,
			Campus:         reqIPAsset.Campus,
			College:        reqIPAsset.College,
			Program:        reqIPAsset.Program,
			Authors:        reqIPAsset.Authors,
		},
	)

	if res.RowsAffected == 0 {
		return nil, errors.New("IP asset not found")
	}

	return &pb.UpdateIP_AssetResponse{
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
		},
	}, nil
}

func (*server) DeleteIPAsset(ctx context.Context, req *pb.DeleteIP_AssetRequest) (*pb.DeleteIP_AssetResponse, error) {
	fmt.Println("Delete IP Asset")
	var ipAsset IP_Asset
	res := DB.Table("table_ipassets").Where("RegistartionNumber = ?", req.GetRegistrationNumber()).Delete(&ipAsset)
	if res.RowsAffected == 0 {
		return nil, errors.New("IP asset not found")
	}

	return &pb.DeleteIP_AssetResponse{
		Success: true,
	}, nil
}

// Publication
func (*server) CreatePublication(ctx context.Context, req *pb.CreatePublicationRequest) (*pb.CreatePublicationResponse, error) {
	fmt.Println("Create Publication")
	publication := req.GetPublication()
	publication.PublicationId = uuid.New().String()

	data := Publication{
		PublicationID:        publication.GetPublicationId(),
		DatePublished:        publication.GetDatePublished(),
		Quartile:             publication.GetQuartile(),
		Authors:              publication.GetAuthors(),
		Department:           publication.GetDepartment(),
		College:              publication.GetCollege(),
		Campus:               publication.GetCampus(),
		TitleOfPaper:         publication.GetTitleOfPaper(),
		TypeOfPublication:    publication.GetTypeOfPublication(),
		FundingSource:        publication.GetFundingSource(),
		NumberOfCitations:    publication.GetNumberOfCitations(),
		GoogleScholarDetails: publication.GetGoogleScholarDetails(),
		SDGNo:                publication.GetSdgNo(),
		FundingType:          publication.GetFundingType(),
		NatureOfFundings:     publication.GetNatureOfFundings(),
		Publisher:            publication.GetPublisher(),
		Abstract:             publication.GetAbstract(),
	}

	res := DB.Table("table_publications").Create(&data)
	if res.RowsAffected == 0 {
		return nil, errors.New("publication creation unsuccessful")
	}

	return &pb.CreatePublicationResponse{
		Publication: &pb.Publication{
			PublicationId:        publication.GetPublicationId(),
			DatePublished:        publication.GetDatePublished(),
			Quartile:             publication.GetQuartile(),
			Authors:              publication.GetAuthors(),
			Department:           publication.GetDepartment(),
			College:              publication.GetCollege(),
			Campus:               publication.GetCampus(),
			TitleOfPaper:         publication.GetTitleOfPaper(),
			TypeOfPublication:    publication.GetTypeOfPublication(),
			FundingSource:        publication.GetFundingSource(),
			NumberOfCitations:    publication.GetNumberOfCitations(),
			GoogleScholarDetails: publication.GetGoogleScholarDetails(),
			SdgNo:                publication.GetSdgNo(),
			FundingType:          publication.GetFundingType(),
			NatureOfFundings:     publication.GetNatureOfFundings(),
			Publisher:            publication.GetPublisher(),
			Abstract:             publication.GetAbstract(),
		},
	}, nil
}

func (*server) GetPublication(ctx context.Context, req *pb.ReadPublicationRequest) (*pb.ReadPublicationResponse, error) {
	fmt.Println("Read Publication", req.GetPublicationId())
	var publication Publication
	res := DB.Table("table_publications").Find(&publication, "PublicationID = ?", req.GetPublicationId())
	if res.RowsAffected == 0 {
		return nil, errors.New("publication not found")
	}

	return &pb.ReadPublicationResponse{
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
			NumberOfCitations:    publication.NumberOfCitations,
			GoogleScholarDetails: publication.GoogleScholarDetails,
			SdgNo:                publication.SDGNo,
			FundingType:          publication.FundingType,
			NatureOfFundings:     publication.NatureOfFundings,
			Publisher:            publication.Publisher,
			Abstract:             publication.Abstract,
		},
	}, nil
}

func (*server) GetPublications(ctx context.Context, req *pb.ReadPublicationsRequest) (*pb.ReadPublicationsResponse, error) {
	fmt.Println("Read Publications")
	publications := []*pb.Publication{}
	res := DB.Table("table_publications").Find(&publications)
	if res.RowsAffected == 0 {
		return nil, errors.New("publications not found")
	}

	return &pb.ReadPublicationsResponse{
		Publications: publications,
	}, nil
}

func (*server) UpdatePublication(ctx context.Context, req *pb.UpdatePublicationRequest) (*pb.UpdatePublicationResponse, error) {
	fmt.Println("Update Publication")
	var publication Publication
	reqPublication := req.GetPublication()

	res := DB.Table("table_publications").Model(&publication).Where("PublicationID = ?", reqPublication.PublicationId).Updates(
		Publication{
			DatePublished:        reqPublication.DatePublished,
			Quartile:             reqPublication.Quartile,
			Authors:              reqPublication.Authors,
			Department:           reqPublication.Department,
			College:              reqPublication.College,
			Campus:               reqPublication.Campus,
			TitleOfPaper:         reqPublication.TitleOfPaper,
			TypeOfPublication:    reqPublication.TypeOfPublication,
			FundingSource:        reqPublication.FundingSource,
			NumberOfCitations:    reqPublication.NumberOfCitations,
			GoogleScholarDetails: reqPublication.GoogleScholarDetails,
			SDGNo:                reqPublication.SdgNo,
			FundingType:          reqPublication.FundingType,
			NatureOfFundings:     reqPublication.NatureOfFundings,
			Publisher:            reqPublication.Publisher,
			Abstract:             reqPublication.Abstract,
		})

	if res.RowsAffected == 0 {
		return nil, errors.New("publication not found")
	}

	return &pb.UpdatePublicationResponse{
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
			NumberOfCitations:    publication.NumberOfCitations,
			GoogleScholarDetails: publication.GoogleScholarDetails,
			SdgNo:                publication.SDGNo,
			FundingType:          publication.FundingType,
			NatureOfFundings:     publication.NatureOfFundings,
			Publisher:            publication.Publisher,
			Abstract:             publication.Abstract,
		},
	}, nil
}

func (*server) DeletePublication(ctx context.Context, req *pb.DeletePublicationRequest) (*pb.DeletePublicationResponse, error) {
	fmt.Println("Delete Publication")
	var publication Publication
	res := DB.Table("table_publications").Where("PublicationID = ?", req.GetPublicationId()).Delete(&publication)
	if res.RowsAffected == 0 {
		return nil, errors.New("publication not found")
	}

	return &pb.DeletePublicationResponse{
		Success: true,
	}, nil
}

// User
func (*server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	fmt.Println("Create User")
	user := req.GetUser()

	data := User{
		UserID:      user.GetUserId(),
		SRCode:      user.GetSrCode(),
		Email:       user.GetEmail(),
		Password:    user.GetPassword(),
		AccountType: user.GetAccountType(),
		UserContact: user.GetUserContact(),
		UserImg:     user.GetUserImg(),
		UserFname:   user.GetUserFname(),
		UserLname:   user.GetUserLname(),
		UserMname:   user.GetUserMname(),
	}

	res := DB.Table("table_user").Create(&data)
	if res.RowsAffected == 0 {
		return nil, errors.New("user creation unsuccessful")
	}

	return &pb.CreateUserResponse{
		User: &pb.User{
			UserId:      user.GetUserId(),
			SrCode:      user.GetSrCode(),
			Email:       user.GetEmail(),
			Password:    user.GetPassword(),
			AccountType: user.GetAccountType(),
			UserContact: user.GetUserContact(),
			UserImg:     user.GetUserImg(),
			UserFname:   user.GetUserFname(),
			UserLname:   user.GetUserLname(),
			UserMname:   user.GetUserMname(),
		},
	}, nil
}

func (*server) GetUser(ctx context.Context, req *pb.ReadUserRequest) (*pb.ReadUserResponse, error) {
	fmt.Println("Read User", req.GetUserId())
	var user User
	res := DB.Table("table_user").Find(&user, "UserID = ?", req.GetUserId())
	if res.RowsAffected == 0 {
		return nil, errors.New("User not found")
	}

	return &pb.ReadUserResponse{
		User: &pb.User{
			UserId:      user.UserID,
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
	}, nil
}

func (*server) GetUsers(ctx context.Context, req *pb.ReadUsersRequest) (*pb.ReadUsersResponse, error) {
	fmt.Println("Read Users")
	users := []*pb.User{}
	res := DB.Table("table_user").Find(&users)
	if res.RowsAffected == 0 {
		return nil, errors.New("Users not found")
	}

	return &pb.ReadUsersResponse{
		Users: users,
	}, nil
}

func (*server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	fmt.Println("Update User")
	var user User
	reqUser := req.GetUser()

	res := DB.Table("table_user").Model(&user).Where("UserID = ?", reqUser.GetUserId()).Updates(
		User{
			SRCode:      reqUser.SrCode,
			Email:       reqUser.Email,
			Password:    reqUser.Password,
			AccountType: reqUser.AccountType,
			UserContact: reqUser.UserContact,
			UserImg:     reqUser.UserImg,
			UserFname:   reqUser.UserFname,
			UserLname:   reqUser.UserLname,
			UserMname:   reqUser.UserMname,
		})
	if res.RowsAffected == 0 {
		return nil, errors.New("User not found")
	}

	return &pb.UpdateUserResponse{
		User: &pb.User{
			UserId:      user.UserID,
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
	}, nil
}

func (*server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	fmt.Println("Delete User")
	var user User
	res := DB.Table("table_user").Where("UserID = ?", req.GetUserId()).Delete(&user)
	if res.RowsAffected == 0 {
		return nil, errors.New("User not found")
	}

	return &pb.DeleteUserResponse{
		Success: true,
	}, nil
}

func main() {
	fmt.Println("gRPC server running ...")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterRMSServiceServer(s, &server{})

	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
}
