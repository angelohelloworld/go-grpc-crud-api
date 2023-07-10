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
	CreatedAt    time.Time `gorm:"autoCreateTime:true"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime:true"`
}

type IP_Asset struct {
	RegistrationNumber string `gorm:"primarykey"`
	TitleOfWork        string
	TypeOfDocument     string
	ClassOfWork        string
	DateOfCreation     string
	DateRegistered     string
	Campus             string
	College            string
	Program            string
	Authors            string
	Hyperlink          string
	Status             string
	Certificate        string
	CreatedAt          time.Time `gorm:"autoCreateTime:true"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime:true"`
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
	NumberOfCitation     int32
	GoogleScholarDetails string
	SDGNo                string
	FundingType          string
	NatureOfFunding      string
	Publisher            string
	Abstract             string
	CreatedAt            time.Time `gorm:"autoCreateTime:true"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime:true"`
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
	CreatedAt   time.Time `gorm:"autoCreateTime:true"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime:true"`
}

type Log struct {
	LogID       string `gorm:"primarykey"`
	DateTime    string
	UserID      int32
	Activity    string
	Description string
	CreatedAt   time.Time `gorm:"autoCreateTime:true"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime:true"`
}

func DatabaseConnection() {
	host := "127.0.0.1"
	port := "5432"
	dbName := "rms_db"
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
	DB.AutoMigrate(&Log{})

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
	res := DB.Table("table_authors").Find(&author, "author_id = ?", req.GetAuthorId())
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

	res := DB.Table("table_authors").Model(&author).Where("author_id = ?", reqAuthor.AuthorId).Updates(
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
	res := DB.Table("table_authors").Where("author_id = ?", req.GetAuthorId()).Delete(&author)
	if res.RowsAffected == 0 {
		return nil, errors.New("author not found")
	}

	return &pb.DeleteAuthorResponse{
		Success: true,
	}, nil
}

// IP_Asset
func (*server) CreateIP_Asset(ctx context.Context, req *pb.CreateIP_AssetRequest) (*pb.CreateIP_AssetResponse, error) {
	fmt.Println("Create IP_Asset")
	ipAsset := req.GetIpAsset()
	ipAsset.RegistrationNumber = uuid.New().String()

	data := IP_Asset{
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
		Hyperlink:          ipAsset.GetHyperlink(),
		Status:             ipAsset.GetStatus(),
		Certificate:        ipAsset.GetCertificate(),
	}

	res := DB.Table("table_ipassets").Create(&data)
	if res.RowsAffected == 0 {
		return nil, errors.New("IP_asset creation unsuccessful")
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
			Hyperlink:          ipAsset.GetHyperlink(),
			Status:             ipAsset.GetStatus(),
			Certificate:        ipAsset.GetCertificate(),
		},
	}, nil
}

func (*server) GetIP_Asset(ctx context.Context, req *pb.ReadIP_AssetRequest) (*pb.ReadIP_AssetResponse, error) {
	fmt.Println("Read IP_assets", req.GetRegistrationNumber())
	var ipAsset IP_Asset
	res := DB.Table("table_ipassets").Find(&ipAsset, "registration_number = ?", req.GetRegistrationNumber())
	if res.RowsAffected == 0 {
		return nil, errors.New("IP_asset not found")
	}
	return &pb.ReadIP_AssetResponse{
		IpAsset: &pb.IP_Asset{
			RegistrationNumber: ipAsset.RegistrationNumber,
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
	}, nil
}

func (*server) GetIP_Assets(ctx context.Context, req *pb.ReadIP_AssetsRequest) (*pb.ReadIP_AssetsResponse, error) {
	fmt.Println("Read IP_assets")
	ipAssets := []*pb.IP_Asset{}
	res := DB.Table("table_ipassets").Find(&ipAssets)
	if res.RowsAffected == 0 {
		return nil, errors.New("IP_asset not found")
	}
	return &pb.ReadIP_AssetsResponse{
		IpAssets: ipAssets,
	}, nil
}

func (*server) UpdateIP_Asset(ctx context.Context, req *pb.UpdateIP_AssetRequest) (*pb.UpdateIP_AssetResponse, error) {
	fmt.Println("Update IP_assets")
	var ipAsset IP_Asset
	reqIPAsset := req.GetIpAsset()

	res := DB.Table("table_ipassets").Model(&ipAsset).Where("registration_number = ?", reqIPAsset.RegistrationNumber).Updates(
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
			Hyperlink:      reqIPAsset.Hyperlink,
			Status:         reqIPAsset.Status,
			Certificate:    reqIPAsset.Certificate,
		})

	if res.RowsAffected == 0 {
		return nil, errors.New("IP_asset not found")
	}

	return &pb.UpdateIP_AssetResponse{
		IpAsset: &pb.IP_Asset{
			RegistrationNumber: ipAsset.RegistrationNumber,
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
	}, nil
}

func (*server) DeleteIP_Asset(ctx context.Context, req *pb.DeleteIP_AssetRequest) (*pb.DeleteIP_AssetResponse, error) {
	fmt.Println("Delete IP_assets")
	var ipAsset IP_Asset
	res := DB.Table("table_ipassets").Where("registration_number = ?", req.GetRegistrationNumber()).Delete(&ipAsset)
	if res.RowsAffected == 0 {
		return nil, errors.New("IP_asset not found")
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
		NumberOfCitation:     publication.GetNumberOfCitation(),
		GoogleScholarDetails: publication.GetGoogleScholarDetails(),
		SDGNo:                publication.GetSdgNo(),
		FundingType:          publication.GetFundingType(),
		NatureOfFunding:      publication.GetNatureOfFunding(),
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
			NumberOfCitation:     publication.GetNumberOfCitation(),
			GoogleScholarDetails: publication.GetGoogleScholarDetails(),
			SdgNo:                publication.GetSdgNo(),
			FundingType:          publication.GetFundingType(),
			NatureOfFunding:      publication.GetNatureOfFunding(),
			Publisher:            publication.GetPublisher(),
			Abstract:             publication.GetAbstract(),
		},
	}, nil
}

func (*server) GetPublication(ctx context.Context, req *pb.ReadPublicationRequest) (*pb.ReadPublicationResponse, error) {
	fmt.Println("Read Publication", req.GetPublicationId())
	var publication Publication
	res := DB.Table("table_publications").Find(&publication, "publication_id = ?", req.GetPublicationId())
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
			NumberOfCitation:     publication.NumberOfCitation,
			GoogleScholarDetails: publication.GoogleScholarDetails,
			SdgNo:                publication.SDGNo,
			FundingType:          publication.FundingType,
			NatureOfFunding:      publication.NatureOfFunding,
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

	res := DB.Table("table_publications").Model(&publication).Where("publication_id = ?", reqPublication.PublicationId).Updates(
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
			NumberOfCitation:     reqPublication.NumberOfCitation,
			GoogleScholarDetails: reqPublication.GoogleScholarDetails,
			SDGNo:                reqPublication.SdgNo,
			FundingType:          reqPublication.FundingType,
			NatureOfFunding:      reqPublication.NatureOfFunding,
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
			NumberOfCitation:     publication.NumberOfCitation,
			GoogleScholarDetails: publication.GoogleScholarDetails,
			SdgNo:                publication.SDGNo,
			FundingType:          publication.FundingType,
			NatureOfFunding:      publication.NatureOfFunding,
			Publisher:            publication.Publisher,
			Abstract:             publication.Abstract,
		},
	}, nil
}

func (*server) DeletePublication(ctx context.Context, req *pb.DeletePublicationRequest) (*pb.DeletePublicationResponse, error) {
	fmt.Println("Delete Publication")
	var publication Publication
	res := DB.Table("table_publications").Where("publication_id = ?", req.GetPublicationId()).Delete(&publication)
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
	res := DB.Table("table_user").Find(&user, "user_id = ?", req.GetUserId())
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

	res := DB.Table("table_user").Model(&user).Where("user_id = ?", reqUser.GetUserId()).Updates(
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
	res := DB.Table("table_user").Where("user_id = ?", req.GetUserId()).Delete(&user)
	if res.RowsAffected == 0 {
		return nil, errors.New("User not found")
	}

	return &pb.DeleteUserResponse{
		Success: true,
	}, nil
}

// Log
func (*server) CreateLog(ctx context.Context, req *pb.CreateLogRequest) (*pb.CreateLogResponse, error) {
	fmt.Println("Create Log")
	log := req.GetLog()

	data := Log{
		LogID:       log.GetLogId(),
		DateTime:    log.GetDateTime(),
		UserID:      log.GetUserId(),
		Activity:    log.GetActivity(),
		Description: log.GetDescription(),
	}

	res := DB.Table("table_log").Create(&data)
	if res.RowsAffected == 0 {
		return nil, errors.New("log creation unsuccessful")
	}

	return &pb.CreateLogResponse{
		Log: &pb.Log{
			LogId:       log.GetLogId(),
			DateTime:    log.GetDateTime(),
			UserId:      log.GetUserId(),
			Activity:    log.GetActivity(),
			Description: log.GetDescription(),
		},
	}, nil
}

func (*server) GetLog(ctx context.Context, req *pb.ReadLogRequest) (*pb.ReadLogResponse, error) {
	fmt.Println("Read Log", req.GetLogId())
	var log Log
	res := DB.Table("table_log").Find(&log, "log_id = ?", req.GetLogId())
	if res.RowsAffected == 0 {
		return nil, errors.New("Log not found")
	}

	return &pb.ReadLogResponse{
		Log: &pb.Log{
			LogId:       log.LogID,
			DateTime:    log.DateTime,
			UserId:      log.UserID,
			Activity:    log.Activity,
			Description: log.Description,
		},
	}, nil
}

func (*server) GetLogs(ctx context.Context, req *pb.ReadLogsRequest) (*pb.ReadLogsResponse, error) {
	fmt.Println("Read Logs")
	logs := []*pb.Log{}
	res := DB.Table("table_log").Find(&logs)
	if res.RowsAffected == 0 {
		return nil, errors.New("Logs not found")
	}

	return &pb.ReadLogsResponse{
		Logs: logs,
	}, nil
}

func (*server) UpdateLog(ctx context.Context, req *pb.UpdateLogRequest) (*pb.UpdateLogResponse, error) {
	fmt.Println("Update Log")
	var log Log
	reqLog := req.GetLog()

	res := DB.Table("table_log").Model(&log).Where("log_id = ?", reqLog.GetLogId()).Updates(
		Log{
			DateTime:    reqLog.DateTime,
			UserID:      reqLog.UserId,
			Activity:    reqLog.Activity,
			Description: reqLog.Description,
		})
	if res.RowsAffected == 0 {
		return nil, errors.New("Log not found")
	}

	return &pb.UpdateLogResponse{
		Log: &pb.Log{
			LogId:       log.LogID,
			DateTime:    log.DateTime,
			UserId:      log.UserID,
			Activity:    log.Activity,
			Description: log.Description,
		},
	}, nil
}

func (*server) DeleteLog(ctx context.Context, req *pb.DeleteLogRequest) (*pb.DeleteLogResponse, error) {
	fmt.Println("Delete Log")
	var log Log
	res := DB.Table("table_log").Where("log_id = ?", req.GetLogId()).Delete(&log)
	if res.RowsAffected == 0 {
		return nil, errors.New("Log not found")
	}

	return &pb.DeleteLogResponse{
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

// deploy server command
// go run server/main.go
// run client command
// go run client/main.go
