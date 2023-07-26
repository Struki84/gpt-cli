package read

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/embeddings/openai"
	llm "github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/pinecone"
)

func Prompt(query string, path string) {

	// >>>>> Load and split PDF
	file, err :=  os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	
	defer file.Close()

	fileInfo, _ := file.Stat()

	PDFLoader := documentloaders.NewPDF(file, fileInfo.Size())

	split := textsplitter.NewRecursiveCharacter()
	split.ChunkSize = 500
	split.ChunkOverlap = 50
	
	docs, err := PDFLoader.LoadAndSplit(context.Background(), split)
	if err != nil {
		fmt.Println(err)
	}

	// >>>>> Embeddings
	e, err := openai.NewOpenAI()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Pinecone vector store.
	store, err := pinecone.New(
		context.Background(),
		pinecone.WithNameSpace(uuid.New().String()),
		pinecone.WithProjectName("fd4e2b9"),
		pinecone.WithAPIKey("65ae7457-8f3d-4b23-a54e-d19b827ab218"),
		pinecone.WithEnvironment("us-west1-gcp-free"),
		pinecone.WithIndexName("reading-test"),
		pinecone.WithEmbedder(e),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = store.AddDocuments(context.Background(), docs)
	if err != nil {
		fmt.Println(err)
	}

	docs, err = store.SimilaritySearch(
		context.Background(), 
		query, 
		5, 
		vectorstores.WithScoreThreshold(0.40),
	)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Similarity search Docs: ", docs)

	// >>>>> Calls 
	llm, err := llm.New()
	if err != nil {
		fmt.Println(err)
	}
	
	stuffQAChain := chains.LoadStuffQA(llm)
	
	answer, err := chains.Call(context.Background(), stuffQAChain, map[string]any{
		"input_documents": docs,
		"question": query,
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(answer["text"])
}