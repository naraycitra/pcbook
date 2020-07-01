package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/naraycitra/pcbook/pb"
	"github.com/naraycitra/pcbook/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const (
	secretKey             = "rahasia"
	tokenDuration         = 15 * time.Minute
	serverCertificateFile = "certs/server.crt"
	serverKeyFile         = "certs/server.key"
	clientCAFile          = "certs/ca.crt"
)

func seedUser(userStore service.UserStore) error {
	err := createUser(userStore, "admin", "admin", "admin")
	if err != nil {
		return err
	}
	err = createUser(userStore, "user", "user", "user")
	if err != nil {
		return err
	}
	return nil
}

func createUser(userStore service.UserStore, username, password, role string) error {
	user, err := service.NewUser(username, password, role)
	if err != nil {
		return err
	}
	return userStore.Save(user)
}

func accessibleRole() map[string][]string {
	const laptopServicePath = "/pcbook.LaptopService/"
	return map[string][]string{
		laptopServicePath + "CreateLaptop": {"admin"},
		laptopServicePath + "UploadImage":  {"admin"},
		laptopServicePath + "RateLaptop":   {"admin", "user"},
	}
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(serverCertificateFile, serverKeyFile)
	if err != nil {
		return nil, err
	}

	// Load certificate of the CA who signed client's certificate
	caCert, err := ioutil.ReadFile(clientCAFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("fail to add client CA's certificate: %v", err)
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}

func runGRPCServer(
	authServer pb.AuthServiceServer,
	laptopServer pb.LaptopServiceServer,
	jwtManager *service.JWTManager,
	enableTLS bool,
	listener net.Listener) error {

	interceptor := service.NewAuthInterceptor(jwtManager, accessibleRole())

	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	}

	if enableTLS {
		tlsCredential, err := loadTLSCredentials()
		if err != nil {
			return fmt.Errorf("cannot load TLS cert: %v", err)
		}
		serverOptions = append(serverOptions, grpc.Creds(tlsCredential))

	}

	grpcServer := grpc.NewServer(serverOptions...)
	pb.RegisterAuthServiceServer(grpcServer, authServer)
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	reflection.Register(grpcServer)

	log.Printf("start GRPC Server on port %s, TLS = %t", listener.Addr().String(), enableTLS)
	return grpcServer.Serve(listener)
}

func runRESTServer(
	authServer pb.AuthServiceServer,
	laptopServer pb.LaptopServiceServer,
	jwtManager *service.JWTManager,
	enableTLS bool,
	listener net.Listener,
	grpcEnpoint string) error {

	mux := runtime.NewServeMux()
	dialOptions := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// in-process handler
	// err := pb.RegisterAuthServiceHandlerServer(ctx, mux, authServer)
	err := pb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, grpcEnpoint, dialOptions)
	if err != nil {
		return err
	}

	// err = pb.RegisterLaptopServiceHandlerServer(ctx, mux, laptopServer)
	err = pb.RegisterLaptopServiceHandlerFromEndpoint(ctx, mux, grpcEnpoint, dialOptions)
	if err != nil {
		return err
	}

	log.Printf("start REST Server on port %s, TLS = %t", listener.Addr().String(), enableTLS)

	if enableTLS {
		return http.ServeTLS(listener, mux, serverCertificateFile, serverKeyFile)
	}
	return http.Serve(listener, mux)
}

func main() {
	port := flag.Int("port", 0, "the server port")
	enableTLS := flag.Bool("tls", false, "enable SSL/TLS")
	serverType := flag.String("type", "grpc", "type of server (grpc/rest)")
	endPoint := flag.String("endpoint", "", "gRPC endpoint")
	flag.Parse()
	log.Printf("start server on port %d, TLS = %t", *port, *enableTLS)

	userStore := service.NewInMemoryUserStore()
	err := seedUser(userStore)
	if err != nil {
		log.Fatalf("cannot seed users")
	}
	jwtManager := service.NewJWTManager(secretKey, tokenDuration)
	authServer := service.NewAuthServer(userStore, jwtManager)

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore("img")
	ratingStore := service.NewInMemoryRatingStore()
	laptopServer := service.NewLaptopServer(laptopStore, imageStore, ratingStore)
	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}

	if *serverType == "grpc" {
		err = runGRPCServer(authServer, laptopServer, jwtManager, *enableTLS, listener)
	} else {
		err = runRESTServer(authServer, laptopServer, jwtManager, *enableTLS, listener, *endPoint)
	}
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}

}
