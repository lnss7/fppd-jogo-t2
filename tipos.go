//aqui ta as informações que os outros usuários e o servidor precisa saber!!!

package main

import "time"

type Jogador struct {
	Nome string
	X    int
	Y    int
	UltimaAtualizacao time.Time
}

type Mensagem struct {
	Remetente string
	Texto     string
}

type Movimento struct {
	Nome           string
	X              int
	Y              int
	SequenceNumber int
}

type EstadoJogo struct {
	Estados []Jogador
}

type Comando struct {
	Nome           string
	X              int
	Y              int
	SequenceNumber int
	Tipo           string // "mover", "interagir"
}

type RespostaComando struct {
	Sucesso bool
	Mensagem string
}

type RegistroJogador struct {
	Nome string
	X    int
	Y    int
}


//