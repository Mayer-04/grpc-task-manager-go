package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	taskpb "github.com/Mayer-04/grpc-task-manager-go/pkg/taskpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TaskClient struct {
	client taskpb.TaskServiceClient
	conn   *grpc.ClientConn
}

func NewTaskClient(address string) (*TaskClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	client := taskpb.NewTaskServiceClient(conn)
	return &TaskClient{
		client: client,
		conn:   conn,
	}, nil
}

func (tc *TaskClient) Close() {
	tc.conn.Close()
}

func main() {
	// Conectar al servidor gRPC
	serverAddr := "localhost:50051"
	if addr := os.Getenv("GRPC_SERVER_ADDR"); addr != "" {
		serverAddr = addr
	}

	client, err := NewTaskClient(serverAddr)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	fmt.Printf("🚀 Conectado al servidor gRPC en %s\n", serverAddr)
	fmt.Println("=== Task Manager Cliente ===")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		showMenu()
		fmt.Print("Selecciona una opción: ")

		if !scanner.Scan() {
			break
		}

		option := strings.TrimSpace(scanner.Text())

		switch option {
		case "1":
			createTaskInteractive(client, scanner)
		case "2":
			getTaskInteractive(client, scanner)
		case "3":
			updateTaskInteractive(client, scanner)
		case "4":
			deleteTaskInteractive(client, scanner)
		case "5":
			markCompleteInteractive(client, scanner)
		case "6":
			listTasksByUserInteractive(client, scanner)
		case "7":
			listAllTasksInteractive(client)
		case "8":
			runDemo(client)
		case "0":
			fmt.Println("👋 ¡Hasta luego!")
			return
		default:
			fmt.Println("❌ Opción inválida. Intenta de nuevo.")
		}

		fmt.Println("\nPresiona Enter para continuar...")
		scanner.Scan()
	}
}

func showMenu() {
	fmt.Println("\n" + strings.Repeat("=", 40))
	fmt.Println("📋 TASK MANAGER - MENÚ PRINCIPAL")
	fmt.Println(strings.Repeat("=", 40))
	fmt.Println("1. ➕ Crear tarea")
	fmt.Println("2. 🔍 Obtener tarea por ID")
	fmt.Println("3. ✏️  Actualizar tarea")
	fmt.Println("4. 🗑️  Eliminar tarea")
	fmt.Println("5. ✅ Marcar tarea como completada")
	fmt.Println("6. 👤 Listar tareas por usuario")
	fmt.Println("7. 📝 Listar todas las tareas")
	fmt.Println("8. 🎯 Demo automático")
	fmt.Println("0. 🚪 Salir")
	fmt.Println(strings.Repeat("=", 40))
}

func createTaskInteractive(client *TaskClient, scanner *bufio.Scanner) {
	fmt.Println("\n➕ CREAR NUEVA TAREA")
	fmt.Println(strings.Repeat("-", 25))

	fmt.Print("👤 User ID: ")
	scanner.Scan()
	userID := strings.TrimSpace(scanner.Text())
	if userID == "" {
		fmt.Println("❌ User ID es requerido")
		return
	}

	fmt.Print("📝 Título: ")
	scanner.Scan()
	title := strings.TrimSpace(scanner.Text())
	if title == "" {
		fmt.Println("❌ Título es requerido")
		return
	}

	fmt.Print("📄 Descripción (opcional): ")
	scanner.Scan()
	description := strings.TrimSpace(scanner.Text())

	fmt.Print("✅ ¿Completada? (y/n, por defecto n): ")
	scanner.Scan()
	completedStr := strings.ToLower(strings.TrimSpace(scanner.Text()))
	completed := completedStr == "y" || completedStr == "yes"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Preparar request
	req := &taskpb.CreateTaskRequest{
		UserId: userID,
		Title:  title,
	}

	if description != "" {
		req.Description = &description
	}
	req.Completed = &completed

	resp, err := client.client.CreateTask(ctx, req)
	if err != nil {
		fmt.Printf("❌ Error creando tarea: %v\n", err)
		return
	}

	fmt.Println("\n✅ ¡Tarea creada exitosamente!")
	printTask(resp.Task)
}

func getTaskInteractive(client *TaskClient, scanner *bufio.Scanner) {
	fmt.Println("\n🔍 OBTENER TAREA")
	fmt.Println(strings.Repeat("-", 20))

	fmt.Print("🆔 Task ID: ")
	scanner.Scan()
	taskID := strings.TrimSpace(scanner.Text())
	if taskID == "" {
		fmt.Println("❌ Task ID es requerido")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.client.GetTask(ctx, &taskpb.GetTaskRequest{Id: taskID})
	if err != nil {
		fmt.Printf("❌ Error obteniendo tarea: %v\n", err)
		return
	}

	fmt.Println("\n📋 Tarea encontrada:")
	printTask(resp.Task)
}

func updateTaskInteractive(client *TaskClient, scanner *bufio.Scanner) {
	fmt.Println("\n✏️ ACTUALIZAR TAREA")
	fmt.Println(strings.Repeat("-", 22))

	fmt.Print("🆔 Task ID: ")
	scanner.Scan()
	taskID := strings.TrimSpace(scanner.Text())
	if taskID == "" {
		fmt.Println("❌ Task ID es requerido")
		return
	}

	fmt.Print("📝 Nuevo título (Enter para mantener actual): ")
	scanner.Scan()
	title := strings.TrimSpace(scanner.Text())

	fmt.Print("📄 Nueva descripción (Enter para mantener actual): ")
	scanner.Scan()
	description := strings.TrimSpace(scanner.Text())

	fmt.Print("✅ Completada? (y/n, Enter para mantener actual): ")
	scanner.Scan()
	completedStr := strings.ToLower(strings.TrimSpace(scanner.Text()))

	req := &taskpb.UpdateTaskRequest{Id: taskID}

	if title != "" {
		req.Title = &title
	}
	if description != "" {
		req.Description = &description
	}
	if completedStr == "y" || completedStr == "yes" {
		completed := true
		req.Completed = &completed
	} else if completedStr == "n" || completedStr == "no" {
		completed := false
		req.Completed = &completed
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.client.UpdateTask(ctx, req)
	if err != nil {
		fmt.Printf("❌ Error actualizando tarea: %v\n", err)
		return
	}

	fmt.Println("\n✅ ¡Tarea actualizada exitosamente!")
	printTask(resp.Task)
}

func deleteTaskInteractive(client *TaskClient, scanner *bufio.Scanner) {
	fmt.Println("\n🗑️ ELIMINAR TAREA")
	fmt.Println(strings.Repeat("-", 20))

	fmt.Print("🆔 Task ID: ")
	scanner.Scan()
	taskID := strings.TrimSpace(scanner.Text())
	if taskID == "" {
		fmt.Println("❌ Task ID es requerido")
		return
	}

	fmt.Print("⚠️  ¿Estás seguro? (y/n): ")
	scanner.Scan()
	confirm := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if confirm != "y" && confirm != "yes" {
		fmt.Println("❌ Operación cancelada")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.client.DeleteTask(ctx, &taskpb.DeleteTaskRequest{Id: taskID})
	if err != nil {
		fmt.Printf("❌ Error eliminando tarea: %v\n", err)
		return
	}

	if resp.Success {
		fmt.Println("✅ Tarea eliminada exitosamente!")
	} else {
		fmt.Printf("❌ Error: %s\n", resp.Message)
	}
}

func markCompleteInteractive(client *TaskClient, scanner *bufio.Scanner) {
	fmt.Println("\n✅ MARCAR COMO COMPLETADA")
	fmt.Println(strings.Repeat("-", 30))

	fmt.Print("🆔 Task ID: ")
	scanner.Scan()
	taskID := strings.TrimSpace(scanner.Text())
	if taskID == "" {
		fmt.Println("❌ Task ID es requerido")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.client.MarkTaskComplete(ctx, &taskpb.MarkTaskCompleteRequest{Id: taskID})
	if err != nil {
		fmt.Printf("❌ Error marcando tarea: %v\n", err)
		return
	}

	fmt.Println("\n✅ ¡Tarea marcada como completada!")
	printTask(resp.Task)
}

func listTasksByUserInteractive(client *TaskClient, scanner *bufio.Scanner) {
	fmt.Println("\n👤 LISTAR TAREAS POR USUARIO")
	fmt.Println(strings.Repeat("-", 35))

	fmt.Print("👤 User ID: ")
	scanner.Scan()
	userID := strings.TrimSpace(scanner.Text())
	if userID == "" {
		fmt.Println("❌ User ID es requerido")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.client.ListTasksByUser(ctx, &taskpb.ListTasksByUserRequest{UserId: userID})
	if err != nil {
		fmt.Printf("❌ Error listando tareas: %v\n", err)
		return
	}

	if len(resp.Tasks) == 0 {
		fmt.Printf("📭 No se encontraron tareas para el usuario %s\n", userID)
		return
	}

	fmt.Printf("\n📋 Tareas del usuario %s (%d encontradas):\n", userID, len(resp.Tasks))
	fmt.Println(strings.Repeat("-", 50))
	for i, task := range resp.Tasks {
		fmt.Printf("\n🔢 Tarea #%d:\n", i+1)
		printTask(task)
	}
}

func listAllTasksInteractive(client *TaskClient) {
	fmt.Println("\n📝 TODAS LAS TAREAS")
	fmt.Println(strings.Repeat("-", 25))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.client.ListAllTasks(ctx, &taskpb.ListAllTasksRequest{})
	if err != nil {
		fmt.Printf("❌ Error listando tareas: %v\n", err)
		return
	}

	if len(resp.Tasks) == 0 {
		fmt.Println("📭 No hay tareas en el sistema")
		return
	}

	fmt.Printf("\n📋 Todas las tareas (%d encontradas):\n", len(resp.Tasks))
	fmt.Println(strings.Repeat("-", 50))
	for i, task := range resp.Tasks {
		fmt.Printf("\n🔢 Tarea #%d:\n", i+1)
		printTask(task)
	}
}

func runDemo(client *TaskClient) {
	fmt.Println("\n🎯 EJECUTANDO DEMO AUTOMÁTICO")
	fmt.Println(strings.Repeat("=", 40))

	ctx := context.Background()
	demoUserID := "demo-user-123"

	// 1. Crear algunas tareas de prueba
	fmt.Println("\n1️⃣ Creando tareas de demo...")

	tasks := []struct {
		title       string
		description string
		completed   bool
	}{
		{"Aprender gRPC", "Estudiar protobuf y implementar servicios", false},
		{"Configurar Docker", "Setup de PostgreSQL con docker-compose", true},
		{"Escribir tests", "Crear tests unitarios e integración", false},
		{"Documentar API", "Crear documentación de la API gRPC", false},
	}

	var createdTaskIDs []string

	for i, taskData := range tasks {
		req := &taskpb.CreateTaskRequest{
			UserId:      demoUserID,
			Title:       taskData.title,
			Description: &taskData.description,
			Completed:   &taskData.completed,
		}

		resp, err := client.client.CreateTask(ctx, req)
		if err != nil {
			fmt.Printf("❌ Error creando tarea %d: %v\n", i+1, err)
			continue
		}

		createdTaskIDs = append(createdTaskIDs, resp.Task.Id)
		fmt.Printf("✅ Creada: %s\n", taskData.title)
		time.Sleep(500 * time.Millisecond) // Pausa para efecto visual
	}

	// 2. Listar tareas del usuario
	fmt.Println("\n2️⃣ Listando tareas del usuario demo...")
	listResp, err := client.client.ListTasksByUser(ctx, &taskpb.ListTasksByUserRequest{
		UserId: demoUserID,
	})
	if err != nil {
		fmt.Printf("❌ Error listando tareas: %v\n", err)
	} else {
		fmt.Printf("📋 Encontradas %d tareas:\n", len(listResp.Tasks))
		for _, task := range listResp.Tasks {
			status := "⏳ Pendiente"
			if task.Completed {
				status = "✅ Completada"
			}
			fmt.Printf("  • %s - %s\n", task.Title, status)
		}
	}

	// 3. Marcar una tarea como completada
	if len(createdTaskIDs) > 0 {
		fmt.Println("\n3️⃣ Marcando primera tarea como completada...")
		markResp, err := client.client.MarkTaskComplete(ctx, &taskpb.MarkTaskCompleteRequest{
			Id: createdTaskIDs[0],
		})
		if err != nil {
			fmt.Printf("❌ Error marcando tarea: %v\n", err)
		} else {
			fmt.Printf("✅ Tarea completada: %s\n", markResp.Task.Title)
		}
	}

	// 4. Actualizar una tarea
	if len(createdTaskIDs) > 1 {
		fmt.Println("\n4️⃣ Actualizando segunda tarea...")
		newTitle := "Aprender gRPC Avanzado"
		newDesc := "Estudiar interceptors, middleware y streaming"

		updateResp, err := client.client.UpdateTask(ctx, &taskpb.UpdateTaskRequest{
			Id:          createdTaskIDs[1],
			Title:       &newTitle,
			Description: &newDesc,
		})
		if err != nil {
			fmt.Printf("❌ Error actualizando tarea: %v\n", err)
		} else {
			fmt.Printf("✅ Tarea actualizada: %s\n", updateResp.Task.Title)
		}
	}

	// 5. Mostrar estado final
	fmt.Println("\n5️⃣ Estado final de las tareas:")
	finalListResp, err := client.client.ListTasksByUser(ctx, &taskpb.ListTasksByUserRequest{
		UserId: demoUserID,
	})
	if err != nil {
		fmt.Printf("❌ Error listando tareas finales: %v\n", err)
	} else {
		for i, task := range finalListResp.Tasks {
			fmt.Printf("\n📋 Tarea #%d:\n", i+1)
			printTask(task)
		}
	}

	fmt.Println("\n🎉 ¡Demo completado!")
}

func printTask(task *taskpb.Task) {
	fmt.Printf("🆔 ID: %s\n", task.Id)
	fmt.Printf("👤 Usuario: %s\n", task.UserId)
	fmt.Printf("📝 Título: %s\n", task.Title)

	if task.Description != nil && *task.Description != "" {
		fmt.Printf("📄 Descripción: %s\n", *task.Description)
	} else {
		fmt.Println("📄 Descripción: (sin descripción)")
	}

	status := "⏳ Pendiente"
	if task.Completed {
		status = "✅ Completada"
	}
	fmt.Printf("📊 Estado: %s\n", status)

	if task.CreatedAt != nil {
		fmt.Printf("📅 Creada: %s\n", task.CreatedAt.AsTime().Format("2006-01-02 15:04:05"))
	}
	if task.UpdatedAt != nil {
		fmt.Printf("🔄 Actualizada: %s\n", task.UpdatedAt.AsTime().Format("2006-01-02 15:04:05"))
	}
}

func readInput(scanner *bufio.Scanner, prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func readBool(scanner *bufio.Scanner, prompt string) bool {
	input := readInput(scanner, prompt)
	return strings.ToLower(input) == "y" || strings.ToLower(input) == "yes"
}

func readInt(scanner *bufio.Scanner, prompt string) (int, error) {
	input := readInput(scanner, prompt)
	return strconv.Atoi(input)
}
