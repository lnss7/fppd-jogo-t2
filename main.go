// main.go - Loop principal do jogo
package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"time"
)

func main() {
	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	// Inicializa o jogo
	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	cliente, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		panic("Erro ao conectar no servidor: " + err.Error())
	}

	interfaceFinalizar()
	var nome string
	fmt.Print("Digite seu nome: ")
	fmt.Scanln(&nome)
	interfaceIniciar()
	jogador := Jogador{Nome: nome, X: jogo.PosX, Y: jogo.PosY}
	var ok bool
	if err := cliente.Call("ServidorJogo.RegistrarJogador", &jogador, &ok); err != nil {
		log.Fatal("Erro ao registrar jogador:", err)
	}
	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	// Loop principal de entrada
	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
			break
		}

		// envia a nova posição para o servidor sempre que houver movimento
		var ok bool
		mov := Movimento{Nome: nome, X: jogo.PosX, Y: jogo.PosY}
		if err := cliente.Call("ServidorJogo.AtualizarPosicao", &mov, &ok); err != nil {
			log.Println("Erro ao atualizar posicao:", err)
		}

		var estado EstadoJogo
		if err := cliente.Call("ServidorJogo.ObterEstado", &nome, &estado); err != nil {
			log.Println("Erro ao obter estado:", err)
		}

		for _, j := range estado.Estados {
			if j.Nome != nome {
				jogo.Mapa[j.Y][j.X] = Personagem
			}
		}
		interfaceDesenharJogo(&jogo)
	}
}
