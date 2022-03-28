package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoramento = 2
const delay = 2        //Delay entre um monitoramento e outro
const divisorDelay = 3 //Delay entre testes é igual delay/divisorDelay

func main() {
	exibeIntrodução()
	fmt.Println("")
	for {
		exibeMenu()
		switch leComando() {
		case 1:
			iniciarMonitoramento()
		case 2:
			fmt.Println("Exibindo logs...")
			imprimeLogs()
		case 0:
			fmt.Println("Saindo do programa")
			os.Exit(0)
		default:
			fmt.Println("Não conheço esse comando!")
			os.Exit(-1)
		}
	}
}

func devolveNomeEIdade() (string, int) {
	nome := "Raul"
	idade := 27
	return nome, idade
}

func exibeIntrodução() {
	nome := "Raul"
	versao := 1.2
	fmt.Println("ola, sr.", nome, "\n Este programa esta na versão", versao)
}

func exibeMenu() {
	fmt.Println(" 1- Iniciar Monitoramento\n", "2- Exibir Logs\n", "0- Sair do Programa")
}

func leComando() int {
	var comandoLido int
	fmt.Scan(&comandoLido)
	fmt.Println("O comando escolhido foi", comandoLido)
	fmt.Println("")
	return comandoLido
}

func iniciarMonitoramento() {
	fmt.Println("Monitorando...")

	sites := leSitesArquivo()
	for i := 0; i < monitoramento; i++ {
		for i, site := range sites {
			fmt.Println("testando site", i, ":", site)
			testaSite(site)
		}
		fmt.Println("============================================")
		time.Sleep(delay * time.Second)
	}
}

func testaSite(site string) {
	resp, err := http.Get(site)

	if err != nil {
		fmt.Println("Ocorreu um erro:,", err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Site:", site, "carregado com sucesso!")
		registraLog(site, true)
	} else {
		fmt.Println("Site:", site, "esta com problemas. Status code:", resp.StatusCode)
		registraLog(site, false)
	}
	fmt.Println("")
	time.Sleep((delay * time.Second) / divisorDelay)
}

func leSitesArquivo() []string {

	var sites []string

	arquivo, err := os.Open("sites.txt")
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}

	leitor := bufio.NewReader(arquivo)
	for {
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)
		sites = append(sites, linha)
		if err == io.EOF {
			break
		}
	}
	arquivo.Close()

	return sites
}

func registraLog(site string, status bool) {
	arquivo, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") +
		" - " + site + " - Online:" + strconv.FormatBool(status) + "\n")
	arquivo.Close()
}

func imprimeLogs() {

	arquivo, err := ioutil.ReadFile("log.txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(arquivo))
}
