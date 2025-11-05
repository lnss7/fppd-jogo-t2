//aqui ta as informações que os outros usuários e o servidor precisa saber!!!

package main

type Jogador struct {
	Nome string
	X    int
	Y    int
}

type Mensagem struct {
	Remetente string
	Texto     string
}

type Movimento struct {
	Nome string
	X    int
	Y    int
}

type EstadoJogo struct {
	Estados []Jogador
}

// adicionado: wrappers de comando com sequence number
type CmdJogador struct {
	SequenceNumber int
	Jogador        Jogador
}

type CmdMovimento struct {
	SequenceNumber int
	Movimento      Movimento
}

type CmdRemover struct {
	SequenceNumber int
	Nome           string
}
