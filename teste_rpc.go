package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"time"
)

func main() {
	// Verifica argumentos
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run teste_rpc.go <nome_do_jogador>")
		os.Exit(1)
	}

	nomeJogador := os.Args[1]

	// Conecta ao servidor
	clienteRPC, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Erro ao conectar ao servidor:", err)
	}
	defer clienteRPC.Close()

	fmt.Printf("Conectado ao servidor! Testando jogador: %s\n", nomeJogador)

	// Testa registro
	registro := &RegistroJogador{
		Nome: nomeJogador,
		X:    5,
		Y:    5,
	}

	var sucesso bool
	err = clienteRPC.Call("ServidorJogo.RegistrarJogador", registro, &sucesso)
	if err != nil {
		log.Fatal("Erro ao registrar:", err)
	}

	if !sucesso {
		log.Fatal("Falha ao registrar jogador")
	}

	fmt.Printf("Jogador %s registrado com sucesso!\n", nomeJogador)

	// Testa ping
	var ping bool
	err = clienteRPC.Call("ServidorJogo.Ping", &nomeJogador, &ping)
	if err != nil {
		log.Fatal("Erro no ping:", err)
	}
	fmt.Printf("Ping: %v\n", ping)

	// Testa obter estado
	var estado EstadoJogo
	err = clienteRPC.Call("ServidorJogo.ObterEstado", &nomeJogador, &estado)
	if err != nil {
		log.Fatal("Erro ao obter estado:", err)
	}

	fmt.Printf("Estado do jogo: %d jogadores\n", len(estado.Estados))
	for _, jogador := range estado.Estados {
		fmt.Printf("  - %s: (%d, %d)\n", jogador.Nome, jogador.X, jogador.Y)
	}

	// Testa comando
	comando := &Comando{
		Nome:           nomeJogador,
		X:              6,
		Y:              6,
		SequenceNumber: 1,
		Tipo:           "mover",
	}

	var resposta RespostaComando
	err = clienteRPC.Call("ServidorJogo.ProcessarComando", comando, &resposta)
	if err != nil {
		log.Fatal("Erro ao processar comando:", err)
	}

	fmt.Printf("Comando processado: %v - %s\n", resposta.Sucesso, resposta.Mensagem)

	// Aguarda um pouco e testa novamente
	time.Sleep(1 * time.Second)

	err = clienteRPC.Call("ServidorJogo.ObterEstado", &nomeJogador, &estado)
	if err != nil {
		log.Fatal("Erro ao obter estado:", err)
	}

	fmt.Printf("Estado após comando: %d jogadores\n", len(estado.Estados))
	for _, jogador := range estado.Estados {
		fmt.Printf("  - %s: (%d, %d)\n", jogador.Nome, jogador.X, jogador.Y)
	}

	// Desregistra
	err = clienteRPC.Call("ServidorJogo.DesregistrarJogador", &nomeJogador, &sucesso)
	if err != nil {
		log.Fatal("Erro ao desregistrar:", err)
	}

	fmt.Printf("Jogador %s desregistrado: %v\n", nomeJogador, sucesso)
	fmt.Println("Teste concluído com sucesso!")
}
