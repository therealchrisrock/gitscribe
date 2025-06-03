package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {
	// Load environment
	bucketName := os.Getenv("FIREBASE_STORAGE_BUCKET")
	credentialsPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")

	if bucketName == "" {
		bucketName = "teammate-5dbc9.appspot.com"
	}
	if credentialsPath == "" {
		credentialsPath = "./firebase-credentials/teammate-5dbc9-firebase-adminsdk-fbsvc-772ba756b7.json"
	}

	fmt.Printf("Testing Firebase Storage Connection:\n")
	fmt.Printf("Bucket: %s\n", bucketName)
	fmt.Printf("Credentials: %s\n", credentialsPath)

	ctx := context.Background()

	// Create storage client
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
	}
	defer client.Close()

	// Test bucket access
	bucket := client.Bucket(bucketName)

	// Try to get bucket attributes
	attrs, err := bucket.Attrs(ctx)
	if err != nil {
		log.Fatalf("Failed to get bucket attributes: %v", err)
	}

	fmt.Printf("✅ Bucket exists!\n")
	fmt.Printf("Bucket info:\n")
	fmt.Printf("  Name: %s\n", attrs.Name)
	fmt.Printf("  Location: %s\n", attrs.Location)
	fmt.Printf("  StorageClass: %s\n", attrs.StorageClass)
	fmt.Printf("  Created: %v\n", attrs.Created)

	// Try to list some objects (first 10)
	fmt.Printf("\nTesting object listing...\n")
	it := bucket.Objects(ctx, &storage.Query{Prefix: ""})
	count := 0
	for count < 10 {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Printf("Error listing objects: %v\n", err)
			break
		}
		fmt.Printf("  Object: %s (size: %d bytes)\n", objAttrs.Name, objAttrs.Size)
		count++
	}

	if count == 0 {
		fmt.Printf("  No objects found (this is normal for a new bucket)\n")
	} else {
		fmt.Printf("  Found %d objects\n", count)
	}

	fmt.Printf("\n✅ Firebase Storage test completed successfully!\n")
}
