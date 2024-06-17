// package main

// import (
// 	"fmt"
// 	"strings"
// 	"sync"
// 	"time"
// )

// type Book struct {
// 	content     string
// 	mu          sync.Mutex
// 	runDuration time.Duration
// }

// func (b *Book) BookWrite(char string) {
// 	b.mu.Lock()
// 	b.content += char
// 	b.mu.Unlock()
// }

// func (b *Book) BookRead() string {
// 	b.mu.Lock()
// 	defer b.mu.Unlock()
// 	if len(b.content) > 0 {
// 		char := b.content[0]
// 		b.content = b.content[1:]
// 		return string(char)
// 	}
// 	return ""
// }

// type Writer struct {
// 	book           *Book
// 	writeInterval  time.Duration
// 	writtenContent string
// }

// func (w *Writer) Write(text string, stop <-chan bool) error {
// 	ticker := time.NewTicker(w.writeInterval)
// 	defer ticker.Stop()

// 	for _, char := range text {
// 		select {
// 		case <-stop:
// 			return nil
// 		case <-ticker.C:
// 			w.book.BookWrite(string(char))
// 			w.writtenContent += string(char)
// 			fmt.Printf("Written: %s\n", w.writtenContent)
// 		}
// 	}
// 	return nil
// }

// type Reader struct {
// 	book         *Book
// 	readInterval time.Duration
// 	readContent  string
// }

// func (r *Reader) Read(stop <-chan bool) error {
// 	ticker := time.NewTicker(r.readInterval)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-stop:
// 			return nil
// 		case <-ticker.C:
// 			char := r.book.BookRead()
// 			if char == "" {
// 				return nil
// 			}
// 			r.readContent += char
// 			fmt.Printf("Read: %s\n", r.readContent)
// 		}
// 	}
// }

// func countFullWords(writtenText string, originalText string) int {
// 	count := 0
// 	for _, word := range strings.Split(originalText, " ") {
// 		if strings.HasPrefix(writtenText, word) {
// 			count++
// 			if len(writtenText) > len(word) {
// 				writtenText = writtenText[len(word)+1:]
// 			}
// 		} else {
// 			break
// 		}
// 	}
// 	return count
// }

// func main() {
// 	book := &Book{
// 		runDuration: 10 * time.Second, // Set the duration of the process here as part of the book configuration
// 	}
// 	writer := Writer{
// 		book:          book,
// 		writeInterval: 1 * time.Second,
// 	}
// 	reader := Reader{
// 		book:         book,
// 		readInterval: 2 * time.Second,
// 	}

// 	var wg sync.WaitGroup
// 	stop := make(chan bool)

// 	text := "cars magazine"

// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		if err := writer.Write(text, stop); err != nil {
// 			panic(err)
// 		}
// 	}()

// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		if err := reader.Read(stop); err != nil {
// 			panic(err)
// 		}
// 	}()

// 	time.Sleep(book.runDuration)
// 	close(stop)
// 	wg.Wait()

//		fmt.Println()
//		fmt.Printf("Total written characters: %d\n", len(writer.writtenContent))
//		fmt.Printf("Total read characters: %d\n", len(reader.readContent))
//		fmt.Printf("Total written words: %d\n", countFullWords(writer.writtenContent, text))
//		fmt.Printf("Total read words: %d\n", countFullWords(reader.readContent, text))
//	}
package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type Book struct {
	content     string
	mu          sync.Mutex
	runDuration time.Duration
}

func (b *Book) BookWrite(char string) {
	b.mu.Lock()
	b.content += char
	b.mu.Unlock()
}

func (b *Book) BookRead() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.content) > 0 {
		char := b.content[0]
		b.content = b.content[1:]
		return string(char)
	}
	return ""
}

type Writer struct {
	book           *Book
	writeInterval  time.Duration
	writtenContent string
}

func (w *Writer) Write(ctx context.Context, text string) error {
	ticker := time.NewTicker(w.writeInterval)
	defer ticker.Stop()

	for _, char := range text {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			w.book.BookWrite(string(char))
			w.writtenContent += string(char)
			fmt.Printf("Written: %s\n", w.writtenContent)
		}
	}
	return nil
}

type Reader struct {
	book         *Book
	readInterval time.Duration
	readContent  string
}

func (r *Reader) Read(ctx context.Context) error {
	ticker := time.NewTicker(r.readInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			char := r.book.BookRead()
			if char == "" {
				return nil
			}
			r.readContent += char
			fmt.Printf("Read: %s\n", r.readContent)
		}
	}
}

func countFullWords(writtenText string, originalText string) int {
	count := 0
	for _, word := range strings.Split(originalText, " ") {
		if strings.HasPrefix(writtenText, word) {
			count++
			if len(writtenText) > len(word) {
				writtenText = writtenText[len(word)+1:]
			}
		} else {
			break
		}
	}
	return count
}

func main() {
	book := &Book{
		runDuration: 10 * time.Second,
	}
	writer := Writer{
		book:          book,
		writeInterval: 1 * time.Second,
	}
	reader := Reader{
		book:         book,
		readInterval: 2 * time.Second,
	}
	const text = "cars magazine"

	ctx, cancel := context.WithTimeout(context.Background(), book.runDuration)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return writer.Write(ctx, text)
	})

	g.Go(func() error {
		return reader.Read(ctx)
	})

	g.Wait()

	fmt.Println()
	fmt.Printf("Total written characters: %d\n", len(writer.writtenContent))
	fmt.Printf("Total read characters: %d\n", len(reader.readContent))
	fmt.Printf("Total written words: %d\n", countFullWords(writer.writtenContent, text))
	fmt.Printf("Total read words: %d\n", countFullWords(reader.readContent, text))
}
