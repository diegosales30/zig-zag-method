// package main

// import (
// 	"bufio"
// 	"crypto/sha256"
// 	"fmt"
// 	"hash"
// 	"math/big"
// 	"math/rand"
// 	"os"
// 	"os/signal"
// 	"strings"
// 	"sync"
// 	"sync/atomic"
// 	"syscall"
// 	"time"

// 	"github.com/btcsuite/btcd/btcutil/base58"
// 	"github.com/decred/dcrd/dcrec/secp256k1/v4"
// 	"golang.org/x/crypto/ripemd160"
// )

// const (
// 	Reset  = "\033[0m"
// 	Green  = "\033[32m"
// 	Yellow = "\033[33m"
// 	Cyan   = "\033[36m"
// 	Red    = "\033[31m"
// )

// var encontrada atomic.Uint32
// var chavesTestadas atomic.Uint64

// type Alvo struct {
// 	EnderecoTarget string
// 	TargetPKH      [20]byte
// 	RangeMin       *big.Int
// 	RangeMax       *big.Int
// 	RangeLen       *big.Int // O tamanho total do labirinto (Max - Min)
// }

// func main() {
// 	fmt.Print("\033[H\033[2J")
// 	fmt.Println(Yellow + "=======================================================")
// 	fmt.Println("   PROTOCOLO COLMEIA V2: ZIG-ZAG ESTOCÁSTICO (SEED)    ")
// 	fmt.Println("=======================================================" + Reset)

// 	reader := bufio.NewScanner(os.Stdin)

// 	fmt.Printf("%s[1] Endereço BTC do Puzzle:%s ", Cyan, Reset)
// 	reader.Scan()
// 	endereco := strings.TrimSpace(reader.Text())

// 	fmt.Printf("%s[2] Range Mínimo (Hex):%s ", Cyan, Reset)
// 	reader.Scan()
// 	minHex := strings.TrimSpace(reader.Text())

// 	fmt.Printf("%s[3] Range Máximo (Hex):%s ", Cyan, Reset)
// 	reader.Scan()
// 	maxHex := strings.TrimSpace(reader.Text())

// 	fmt.Printf("\n%sDeseja injetar as Formigas no Labirinto de Caos? (Y/N):%s ", Yellow, Reset)
// 	reader.Scan()
// 	if strings.ToLower(strings.TrimSpace(reader.Text())) != "y" {
// 		return
// 	}

// 	decoded, _, err := base58.CheckDecode(endereco)
// 	if err != nil || len(decoded) != 20 {
// 		fmt.Println(Red + "[ERRO] Endereço Bitcoin inválido ou não suportado!" + Reset)
// 		return
// 	}
// 	var targetPKH [20]byte
// 	copy(targetPKH[:], decoded)

// 	minKey := new(big.Int)
// 	maxKey := new(big.Int)
// 	minKey.SetString(minHex, 16)
// 	maxKey.SetString(maxHex, 16)

// 	// Calcula o tamanho real do terreno matemático
// 	rangeLen := new(big.Int).Sub(maxKey, minKey)

// 	alvo := Alvo{
// 		EnderecoTarget: endereco,
// 		TargetPKH:      targetPKH,
// 		RangeMin:       minKey,
// 		RangeMax:       maxKey,
// 		RangeLen:       rangeLen,
// 	}

// 	c := make(chan os.Signal, 1)
// 	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
// 	go func() {
// 		<-c
// 		encontrada.Store(1)
// 		fmt.Println("\n\n" + Red + "[!] RETIRADA TÁTICA ACIONADA." + Reset)
// 		os.Exit(0)
// 	}()

// 	startTime := time.Now()
// 	var wg sync.WaitGroup

// 	go painelMetricas(startTime)

// 	// Lançamento com Sementes independentes para garantir que não sigam a mesma trilha
// 	sementeMestre := time.Now().UnixNano()

// 	wg.Add(1)
// 	go formigaCaosAbsoluto(alvo, &wg, sementeMestre+1)

// 	wg.Add(1)
// 	go formigaZigZagCrescente(alvo, &wg, sementeMestre+2)

// 	wg.Add(1)
// 	go formigaZigZagDecrescente(alvo, &wg, sementeMestre+3)

// 	wg.Add(1)
// 	go formigaOndaCruzada(alvo, &wg, sementeMestre+4)

// 	wg.Wait()

// 	if encontrada.Load() == 0 {
// 		fmt.Printf("\n%s[INFO] Labirinto infinito interrompido.%s\n", Yellow, Reset)
// 	}
// }

// func painelMetricas(start time.Time) {
// 	for encontrada.Load() == 0 {
// 		time.Sleep(1 * time.Second)
// 		elapsed := time.Since(start).Seconds()
// 		chaves := chavesTestadas.Load()
// 		speed := float64(chaves) / elapsed

// 		fmt.Printf("\r%sPoder Estocástico:%s %.0f chaves/s | %sTempo:%s %.1fs ",
// 			Yellow, Reset, speed, Cyan, Reset, elapsed)
// 	}
// }

// // ---------------------------------------------------------------------
// // NÚCLEO DE MATEMÁTICA BARE-METAL
// // ---------------------------------------------------------------------

// func verificarColisaoReal(privKey *big.Int, targetPKH [20]byte, rmd160 hash.Hash) bool {
// 	privBytes := privKey.Bytes()
// 	var padded [32]byte
// 	copy(padded[32-len(privBytes):], privBytes)

// 	priv := secp256k1.PrivKeyFromBytes(padded[:])
// 	pub := priv.PubKey()
// 	pubBytes := pub.SerializeCompressed()

// 	sha := sha256.Sum256(pubBytes)
// 	rmd160.Reset()
// 	rmd160.Write(sha[:])
// 	pkh := rmd160.Sum(nil)

// 	for i := 0; i < 20; i++ {
// 		if pkh[i] != targetPKH[i] {
// 			return false
// 		}
// 	}
// 	return true
// }

// func salvarChaveEncontrada(privKey *big.Int, endereco string) {
// 	encontrada.Store(1) 
// 	privHex := fmt.Sprintf("%064x", privKey)
// 	fmt.Printf("\n\n%s╔══════════════════════════════════════════════════════════╗%s\n", Green, Reset)
// 	fmt.Printf("%s║ RAINHA CAPTURADA NO CAOS! (COLISÃO HASH160)              ║%s\n", Green, Reset)
// 	fmt.Printf("%s║ ENDEREÇO: %-46s ║%s\n", Green, endereco, Reset)
// 	fmt.Printf("%s║ CHAVE: %-49s ║%s\n", Yellow, privHex, Reset)
// 	fmt.Printf("%s╚══════════════════════════════════════════════════════════╝%s\n", Green, Reset)
	
// 	f, err := os.OpenFile("caminho_formigo_v2.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err == nil {
// 		f.WriteString("Alvo: " + endereco + "\nPriv: " + privHex + "\n---------------------\n")
// 		f.Close()
// 	}
// }

// // ---------------------------------------------------------------------
// // AS 4 ROTAS DE CAOS (ZIG-ZAG E RANDOM)
// // ---------------------------------------------------------------------

// // Formiga 0: Pulos 100% caóticos pelo mapa inteiro a cada ciclo.
// func formigaCaosAbsoluto(alvo Alvo, wg *sync.WaitGroup, seed int64) {
// 	defer wg.Done()
// 	rng := rand.New(rand.NewSource(seed))
// 	rmd := ripemd160.New()
	
// 	atual := new(big.Int)
// 	offset := new(big.Int)

// 	for encontrada.Load() == 0 {
// 		// Gera um offset aleatório entre 0 e RangeLen e soma ao RangeMin
// 		offset.Rand(rng, alvo.RangeLen)
// 		atual.Add(alvo.RangeMin, offset)

// 		if verificarColisaoReal(atual, alvo.TargetPKH, rmd) {
// 			salvarChaveEncontrada(atual, alvo.EnderecoTarget)
// 			return
// 		}
// 		chavesTestadas.Add(1)
// 	}
// }

// // Formiga 1: Zig-Zag Crescente (Nasce no caos, anda pra frente com saltos irregulares)
// func formigaZigZagCrescente(alvo Alvo, wg *sync.WaitGroup, seed int64) {
// 	defer wg.Done()
// 	rng := rand.New(rand.NewSource(seed))
// 	rmd := ripemd160.New()
	
// 	atual := new(big.Int)
// 	offset := new(big.Int)
// 	passo := new(big.Int)

// 	// Pinball: Nasce em um lugar aleatório
// 	offset.Rand(rng, alvo.RangeLen)
// 	atual.Add(alvo.RangeMin, offset)

// 	for encontrada.Load() == 0 {
// 		if verificarColisaoReal(atual, alvo.TargetPKH, rmd) {
// 			salvarChaveEncontrada(atual, alvo.EnderecoTarget)
// 			return
// 		}

// 		// Passo embriagado: Pula de 1 a 50.000 posições pra frente
// 		salto := rng.Int63n(50000) + 1
// 		passo.SetInt64(salto)
// 		atual.Add(atual, passo)

// 		// Se bater na parede final (RangeMax), sofre um ricochete para um novo ponto aleatório
// 		if atual.Cmp(alvo.RangeMax) > 0 {
// 			offset.Rand(rng, alvo.RangeLen)
// 			atual.Add(alvo.RangeMin, offset)
// 		}

// 		chavesTestadas.Add(1)
// 	}
// }

// // Formiga 2: Zig-Zag Decrescente (Nasce no caos, recua com saltos irregulares)
// func formigaZigZagDecrescente(alvo Alvo, wg *sync.WaitGroup, seed int64) {
// 	defer wg.Done()
// 	rng := rand.New(rand.NewSource(seed))
// 	rmd := ripemd160.New()
	
// 	atual := new(big.Int)
// 	offset := new(big.Int)
// 	passo := new(big.Int)

// 	offset.Rand(rng, alvo.RangeLen)
// 	atual.Add(alvo.RangeMin, offset)

// 	for encontrada.Load() == 0 {
// 		if verificarColisaoReal(atual, alvo.TargetPKH, rmd) {
// 			salvarChaveEncontrada(atual, alvo.EnderecoTarget)
// 			return
// 		}

// 		// Passo embriagado reverso
// 		salto := rng.Int63n(50000) + 1
// 		passo.SetInt64(salto)
// 		atual.Sub(atual, passo)

// 		// Se bater na parede inicial (RangeMin), sofre um ricochete aleatório
// 		if atual.Cmp(alvo.RangeMin) < 0 {
// 			offset.Rand(rng, alvo.RangeLen)
// 			atual.Add(alvo.RangeMin, offset)
// 		}

// 		chavesTestadas.Add(1)
// 	}
// }

// // Formiga 3: Onda Cruzada (Saltos gigantescos cobrindo frações maciças do terreno)
// func formigaOndaCruzada(alvo Alvo, wg *sync.WaitGroup, seed int64) {
// 	defer wg.Done()
// 	rng := rand.New(rand.NewSource(seed))
// 	rmd := ripemd160.New()
	
// 	atual := new(big.Int)
// 	offset := new(big.Int)
	
// 	// Salto quântico: O tamanho do salto baseia-se em 1% a 5% do mapa total a cada passada
// 	// Como big.Int não divide bem frações menores no Rand, usaremos múltiplos saltos de tamanho randômico absoluto
	
// 	for encontrada.Load() == 0 {
// 		// Calcula um salto massivo no meio do caos para manter o calor distribuído
// 		offset.Rand(rng, alvo.RangeLen)
// 		atual.Add(alvo.RangeMin, offset)

// 		if verificarColisaoReal(atual, alvo.TargetPKH, rmd) {
// 			salvarChaveEncontrada(atual, alvo.EnderecoTarget)
// 			return
// 		}
// 		chavesTestadas.Add(1)
// 	}
// }


// As 4 táticas 
// (Caos Absoluto, Zig-Zag Crescente, Zig-Zag Decrescente e Onda Cruzada)
// serão distribuídas como cartas em uma mesa (Round-Robin) para o número de núcleos que você escolher.
package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"hash"
	"math/big"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"golang.org/x/crypto/ripemd160"
)

const (
	Reset  = "\033[0m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Red    = "\033[31m"
)

var encontrada atomic.Uint32
var chavesTestadas atomic.Uint64

type Alvo struct {
	EnderecoTarget string
	TargetPKH      [20]byte
	RangeMin       *big.Int
	RangeMax       *big.Int
	RangeLen       *big.Int // O tamanho total do labirinto (Max - Min)
}

// Assinatura padrão para as rotas das formigas
type TaticaFormiga func(alvo Alvo, wg *sync.WaitGroup, seed int64, formigaID int)

func main() {
	fmt.Print("\033[H\033[2J")
	fmt.Println(Yellow + "=======================================================")
	fmt.Println("   PROTOCOLO COLMEIA V2: ZIG-ZAG ESTOCÁSTICO MULTICORE ")
	fmt.Println("=======================================================" + Reset)

	// 1. Detecção de Hardware no Metal Nu
	numCores := runtime.NumCPU()
	fmt.Printf("%s[HARDWARE DETECTADO]%s %d Núcleos/Threads de Processamento\n\n", Cyan, Reset, numCores)

	reader := bufio.NewScanner(os.Stdin)

	fmt.Printf("%s[1] Quantos núcleos deseja engajar? (1 a %d):%s ", Cyan, numCores, Reset)
	reader.Scan()
	threadCount, err := strconv.Atoi(strings.TrimSpace(reader.Text()))
	if err != nil || threadCount < 1 {
		threadCount = 1 // Proteção contra entradas inválidas
	}

	fmt.Printf("%s[2] Endereço BTC do Puzzle:%s ", Cyan, Reset)
	reader.Scan()
	endereco := strings.TrimSpace(reader.Text())

	fmt.Printf("%s[3] Range Mínimo (Hex):%s ", Cyan, Reset)
	reader.Scan()
	minHex := strings.TrimSpace(reader.Text())

	fmt.Printf("%s[4] Range Máximo (Hex):%s ", Cyan, Reset)
	reader.Scan()
	maxHex := strings.TrimSpace(reader.Text())

	fmt.Printf("\n%sDeseja injetar as Formigas no Labirinto de Caos? (Y/N):%s ", Yellow, Reset)
	reader.Scan()
	if strings.ToLower(strings.TrimSpace(reader.Text())) != "y" {
		return
	}

	decoded, _, err := base58.CheckDecode(endereco)
	if err != nil || len(decoded) != 20 {
		fmt.Println(Red + "[ERRO] Endereço Bitcoin inválido ou não suportado!" + Reset)
		return
	}
	var targetPKH [20]byte
	copy(targetPKH[:], decoded)

	minKey := new(big.Int)
	maxKey := new(big.Int)
	minKey.SetString(minHex, 16)
	maxKey.SetString(maxHex, 16)

	// Calcula o tamanho real do terreno matemático
	rangeLen := new(big.Int).Sub(maxKey, minKey)

	alvo := Alvo{
		EnderecoTarget: endereco,
		TargetPKH:      targetPKH,
		RangeMin:       minKey,
		RangeMax:       maxKey,
		RangeLen:       rangeLen,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		encontrada.Store(1)
		fmt.Println("\n\n" + Red + "[!] RETIRADA TÁTICA ACIONADA." + Reset)
		os.Exit(0)
	}()

	startTime := time.Now()
	var wg sync.WaitGroup

	go painelMetricas(startTime)

	// Arsenal de Táticas Estocásticas
	taticas := []TaticaFormiga{
		formigaCaosAbsoluto,
		formigaZigZagCrescente,
		formigaZigZagDecrescente,
		formigaOndaCruzada,
	}

	// 2. Lançamento com Sementes independentes (Timestamp Mestre)
	sementeMestre := time.Now().UnixNano()

	// 3. Distribuição das Threads (Round-Robin)
	for i := 0; i < threadCount; i++ {
		wg.Add(1)
		taticaSelecionada := taticas[i%len(taticas)] // Distribui as táticas uniformemente
		sementeIndividual := sementeMestre + int64(i*1000) // Semente temporal exclusiva
		
		go taticaSelecionada(alvo, &wg, sementeIndividual, i)
	}

	wg.Wait()

	if encontrada.Load() == 0 {
		fmt.Printf("\n%s[INFO] Labirinto infinito interrompido.%s\n", Yellow, Reset)
	}
}

func painelMetricas(start time.Time) {
	for encontrada.Load() == 0 {
		time.Sleep(1 * time.Second)
		if encontrada.Load() != 0 { break }
		elapsed := time.Since(start).Seconds()
		chaves := chavesTestadas.Load()
		speed := float64(chaves) / elapsed

		fmt.Printf("\r%sPoder Estocástico:%s %.0f chaves/s | %sTempo:%s %.1fs ",
			Yellow, Reset, speed, Cyan, Reset, elapsed)
	}
}

// ---------------------------------------------------------------------
// NÚCLEO DE MATEMÁTICA BARE-METAL
// ---------------------------------------------------------------------

func verificarColisaoReal(privKey *big.Int, targetPKH [20]byte, rmd160 hash.Hash) bool {
	privBytes := privKey.Bytes()
	var padded [32]byte
	copy(padded[32-len(privBytes):], privBytes)

	priv := secp256k1.PrivKeyFromBytes(padded[:])
	pub := priv.PubKey()
	pubBytes := pub.SerializeCompressed()

	sha := sha256.Sum256(pubBytes)
	rmd160.Reset()
	rmd160.Write(sha[:])
	pkh := rmd160.Sum(nil)

	for i := 0; i < 20; i++ {
		if pkh[i] != targetPKH[i] {
			return false
		}
	}
	return true
}

func salvarChaveEncontrada(privKey *big.Int, endereco string, formigaID int) {
	encontrada.Store(1) 
	privHex := fmt.Sprintf("%064x", privKey)
	fmt.Printf("\n\n%s╔══════════════════════════════════════════════════════════╗%s\n", Green, Reset)
	fmt.Printf("%s║ RAINHA CAPTURADA NO CAOS PELA FORMIGA [%d]!              ║%s\n", Green, formigaID, Reset)
	fmt.Printf("%s║ ENDEREÇO: %-46s ║%s\n", Green, endereco, Reset)
	fmt.Printf("%s║ CHAVE: %-49s ║%s\n", Yellow, privHex, Reset)
	fmt.Printf("%s╚══════════════════════════════════════════════════════════╝%s\n", Green, Reset)
	
	f, err := os.OpenFile("caminho_formigo_v2.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		f.WriteString(fmt.Sprintf("Alvo: %s\nPriv: %s\nFormiga Batedora: %d\n---------------------\n", endereco, privHex, formigaID))
		f.Close()
	}
}

// ---------------------------------------------------------------------
// AS 4 ROTAS DE CAOS (ZIG-ZAG E RANDOM)
// ---------------------------------------------------------------------

// Formiga 0: Pulos 100% caóticos pelo mapa inteiro a cada ciclo.
func formigaCaosAbsoluto(alvo Alvo, wg *sync.WaitGroup, seed int64, id int) {
	defer wg.Done()
	rng := rand.New(rand.NewSource(seed))
	rmd := ripemd160.New()
	
	atual := new(big.Int)
	offset := new(big.Int)

	for encontrada.Load() == 0 {
		// Gera um offset aleatório entre 0 e RangeLen e soma ao RangeMin
		offset.Rand(rng, alvo.RangeLen)
		atual.Add(alvo.RangeMin, offset)

		if verificarColisaoReal(atual, alvo.TargetPKH, rmd) {
			salvarChaveEncontrada(atual, alvo.EnderecoTarget, id)
			return
		}
		chavesTestadas.Add(1)
	}
}

// Formiga 1: Zig-Zag Crescente (Nasce no caos, anda pra frente com saltos irregulares)
func formigaZigZagCrescente(alvo Alvo, wg *sync.WaitGroup, seed int64, id int) {
	defer wg.Done()
	rng := rand.New(rand.NewSource(seed))
	rmd := ripemd160.New()
	
	atual := new(big.Int)
	offset := new(big.Int)
	passo := new(big.Int)

	// Pinball: Nasce em um lugar aleatório
	offset.Rand(rng, alvo.RangeLen)
	atual.Add(alvo.RangeMin, offset)

	for encontrada.Load() == 0 {
		if verificarColisaoReal(atual, alvo.TargetPKH, rmd) {
			salvarChaveEncontrada(atual, alvo.EnderecoTarget, id)
			return
		}

		// Passo embriagado: Pula de 1 a 50.000 posições pra frente
		salto := rng.Int63n(50000) + 1
		passo.SetInt64(salto)
		atual.Add(atual, passo)

		// Se bater na parede final (RangeMax), sofre um ricochete para um novo ponto aleatório
		if atual.Cmp(alvo.RangeMax) > 0 {
			offset.Rand(rng, alvo.RangeLen)
			atual.Add(alvo.RangeMin, offset)
		}

		chavesTestadas.Add(1)
	}
}

// Formiga 2: Zig-Zag Decrescente (Nasce no caos, recua com saltos irregulares)
func formigaZigZagDecrescente(alvo Alvo, wg *sync.WaitGroup, seed int64, id int) {
	defer wg.Done()
	rng := rand.New(rand.NewSource(seed))
	rmd := ripemd160.New()
	
	atual := new(big.Int)
	offset := new(big.Int)
	passo := new(big.Int)

	offset.Rand(rng, alvo.RangeLen)
	atual.Add(alvo.RangeMin, offset)

	for encontrada.Load() == 0 {
		if verificarColisaoReal(atual, alvo.TargetPKH, rmd) {
			salvarChaveEncontrada(atual, alvo.EnderecoTarget, id)
			return
		}

		// Passo embriagado reverso
		salto := rng.Int63n(50000) + 1
		passo.SetInt64(salto)
		atual.Sub(atual, passo)

		// Se bater na parede inicial (RangeMin), sofre um ricochete aleatório
		if atual.Cmp(alvo.RangeMin) < 0 {
			offset.Rand(rng, alvo.RangeLen)
			atual.Add(alvo.RangeMin, offset)
		}

		chavesTestadas.Add(1)
	}
}

// Formiga 3: Onda Cruzada (Saltos gigantescos cobrindo frações maciças do terreno)
func formigaOndaCruzada(alvo Alvo, wg *sync.WaitGroup, seed int64, id int) {
	defer wg.Done()
	rng := rand.New(rand.NewSource(seed))
	rmd := ripemd160.New()
	
	atual := new(big.Int)
	offset := new(big.Int)
	
	for encontrada.Load() == 0 {
		// Calcula um salto massivo no meio do caos para manter o calor distribuído
		offset.Rand(rng, alvo.RangeLen)
		atual.Add(alvo.RangeMin, offset)

		if verificarColisaoReal(atual, alvo.TargetPKH, rmd) {
			salvarChaveEncontrada(atual, alvo.EnderecoTarget, id)
			return
		}
		chavesTestadas.Add(1)
	}
}