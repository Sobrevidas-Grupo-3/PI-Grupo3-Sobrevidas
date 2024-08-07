
package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/jung-kurt/gofpdf"
	_ "github.com/lib/pq"
)

type Pacientes struct {
	Nome        string
	Rua         string
	Numero      string
	Sexo        string
	Bairro      string
	Complemento string
	Telefone    string
	DataNasc    string
	Homem       string
	Etilista    string
	Tabagista   string
	LesaoBucal  string
	Endereco    string
	Idade       int
	Fatores     string
	Usuario     string
	Primeira    string
	TemDados    bool
}

type PacienteFormularioPreenchido struct {
	ID            string
	PrimeiroNome  string
	Nome          string
	DataNasc      string
	CPF           string
	NomeMae       string
	Sexo          string
	CartaoSus     string
	Telefone      string
	Email         string
	CEP           string
	Bairro        string
	Rua           string
	Numero        string
	Complemento   string
	Homem         string
	Etilista      string
	Tabagista     string
	LesaoBucal    string
	Usuario       string
	PrimeiraLetra string
	DataCadastro  string
	UltimaVisita  string
	CNS           []string
	CBO1          string
	CBO2          string
	CBO3          string
	CBO4          string
	CBO5          string
	CBO6          string
	CNES          []string
	INE           []string
	BaixoRisco    bool
	MedioRisco    bool
	AltoRisco     bool
	IsHomem       bool
	IsEtilista    bool
	IsTabagista   bool
	IsLesaoBucal  bool
	MaisDeUmMes   bool
}

type validarlogin struct {
	Usuario         string
	Cpf             string
	Senha           string
	PrimeiraLetra   string
	QtdMaisDeUmMes  int
	QtdTotal        int
	QtdBaixo        int
	QtdMedio        int
	QtdAlto         int
	PorcBaixo       float64
	PorcMedio       float64
	PorcAlto        float64
	Cns             string
	Cbo             string
	Cnes            string
	Ine             string
	Nome            string
	DataNascimento  string
	Etilista        string
	Homem           string
	Telefone        string
	Tabagista       string
	FeridasBucais   string
	Fatores         string
	Baixo           bool
	Medio           bool
	Alto            bool
	UltimaVisita    string
	Cep             string
	Bairro          string
	Rua             string
	Numero          string
	Complemento     string
	Endereco        string
	DadosGoogleMaps []DadosGoogleMaps
}

type DadosGoogleMaps struct {
	Cep            string
	Nome           string
	DataNascimento string
	Telefone       string
	Fatores        string
	Endereco       string
	UltimaVisita   string
	MaisDeUmMes    bool
	Alto           bool
	Medio          bool
	Baixo          bool
}

type PegarDados struct {
	Homem        string
	Etilista     string
	Tabagista    string
	LesaoBucal   string
	UltimaVisita string
}

type UsuarioNoDashboard struct {
	Usuario   string
	Primeira  string
	QtdBaixo  int
	QtdMedio  int
	QtdAlto   int
	PorcBaixo float64
	PorcMedio float64
	PorcAlto  float64
}

type validarCpf struct {
	Cpf string
}

type ACS struct {
	User          string
	NomeCompleto  string
	CPF           string
	CNS           string
	CBO           string
	CNES          string
	INE           string
	SenhaACS      string
	PrimeiraLetra string
}

type DadosForm struct {
	Usuario       string
	PrimeiraLetra string
	Cns           []string
	Cbo1          string
	Cbo2          string
	Cbo3          string
	Cbo4          string
	Cbo5          string
	Cbo6          string
	Cnes          []string
	Ine           []string
	ErroCampos    bool
	Confirmacao   bool
}

var db = fazConexaoComBanco()
var Cns, Cbo, Cnes, Ine, cpfLogin, senhaLogin, usuarioLogin, primeiraletraLogin, nomePaciente string
var loginInvalido = false
var esqueceuInvalido = false
var confirmCadastro = false
var erroCadastro bool
var qtdBaixo, qtdMedio, qtdAlto, qtdTotal, qtdTotalCard, qtdVisitadosMaisDeUmMes int
var templates = template.Must(template.ParseFiles("./index.html", "./templates/cadastro/cadastro.html", "./templates/telalogin/login.html", "./templates/telaesqueceusenha/esqueceusenha.html", "./templates/dashboard/dashboard.html", "./templates/formulario/formulario.html", "./templates/central-usuario/centralusuario.html", "./templates/pacientesgerais/pacGerais.html", "./templates/pg-baixo/pg-baixo.html", "./templates/pg-medio/pg-medio.html", "./templates/pg-alto/pg-alto.html", "./templates/pg-absenteista/pg-absenteista.html", "./templates/pag-Faq/faq.html", "./templates/formulario-preenchido/formpreenchido.html"))

func main() {
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)
	colocarDados()
	http.HandleFunc("/cadastro", executarCadastro)
	http.HandleFunc("/login", autenticaCadastroELevaAoLogin)
	http.HandleFunc("/login/autenticar", loginInvalidado)
	http.HandleFunc("/dashboard", autenticaLoginELevaAoDashboard)
	http.HandleFunc("/dashboard/voltar", dashboard)
	http.HandleFunc("/esqueceusenha", executarEsqueceuSenha)
	http.HandleFunc("/esqueceusenha/invalido", esqueceuSenhaInvalidado)
	http.HandleFunc("/telalogin", atualizarSenha)
	http.HandleFunc("/cadastrar-paciente", executarFormulario)
	http.HandleFunc("/paciente-cadastrado", cadastrarPaciente)
	http.HandleFunc("/central-usuario", executarCentralUsuario)
	http.HandleFunc("/central-usuario/atualizarsenha", atualizarSenhaCentralUsuario)
	http.HandleFunc("/pagina-faq", executarPagFaq)
	http.HandleFunc("/pacientesgerais", executarPacGerais)
	http.HandleFunc("/pagina-baixo-risco", executarPgBaixo)
	http.HandleFunc("/pagina-baixo-risco/filtrar", executarPgBaixoFiltro)
	http.HandleFunc("/pagina-baixo-risco/filtrar-nome", executarPgBaixoFiltroPorNome)
	http.HandleFunc("/pagina-medio-risco", executarPgMedio)
	http.HandleFunc("/pagina-medio-risco/filtrar", executarPgMedioFiltro)
	http.HandleFunc("/pagina-medio-risco/filtrar-nome", executarPgMedioFiltroPorNome)
	http.HandleFunc("/pagina-alto-risco", executarPgAlto)
	http.HandleFunc("/pagina-alto-risco/filtrar", executarPgAltoFiltro)
	http.HandleFunc("/pagina-alto-risco/filtrar-nome", executarPgAltoFiltroPorNome)
	http.HandleFunc("/pagina-absenteista", executarPgAbsenteista)
	http.HandleFunc("/formulario/preenchido", executarFormPreenchido)
	http.HandleFunc("/formulario/preenchido/mapa", executarFormPreenchidoVindoDoMaps)
	http.HandleFunc("/formulario/preenchido/ultimavisita-alterada", alterarDataUltimaVisitaFormPreenchido)
	http.HandleFunc("/generate-pdf", generatePDF)

	log.Println("Server rodando na porta 8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func fazConexaoComBanco() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar arquivo .env")
	}

	usuarioBancoDeDados := os.Getenv("USUARIO")
	senhaDoUsuario := os.Getenv("SENHA")
	nomeDoBancoDeDados := os.Getenv("NOME_BANCO_DE_DADOS")
	dadosParaConexao := "user=" + usuarioBancoDeDados + " dbname=" + nomeDoBancoDeDados + " password=" + senhaDoUsuario + " host=localhost port=5432 sslmode=disable"
	database, err := sql.Open("postgres", dadosParaConexao)
	if err != nil {
		fmt.Println(err)
	}

	_, err = database.Query("CREATE TABLE IF NOT EXISTS cadastro(id SERIAL PRIMARY KEY, nome_completo VARCHAR(255) NOT NULL, cpf VARCHAR(15) UNIQUE NOT NULL, cns VARCHAR(15), cbo VARCHAR(15), cnes VARCHAR(15), ine VARCHAR(15), senha VARCHAR(20))")
	if err != nil {
		log.Fatal(err)
	}

	_, err = database.Query("CREATE TABLE IF NOT EXISTS pacientes(id SERIAL PRIMARY KEY, nome_completo VARCHAR(255), data_nasc VARCHAR(30), cpf VARCHAR(15) UNIQUE NOT NULL, nome_mae VARCHAR(255), sexo VARCHAR(30), cartao_sus VARCHAR(55) UNIQUE NOT NULL, telefone VARCHAR(55) UNIQUE NOT NULL, email VARCHAR(255) UNIQUE NOT NULL, cep VARCHAR(15), bairro VARCHAR(255), rua VARCHAR(255), numero VARCHAR(255), complemento VARCHAR(255), homem VARCHAR(15) NOT NULL, etilista VARCHAR(15) NOT NULL, tabagista VARCHAR(15) NOT NULL, lesao_bucal VARCHAR(15) NOT NULL, data_cadastro VARCHAR(20), ultima_visita VARCHAR(20))")
	if err != nil {
		log.Fatal(err)
	}

	return database
}

func executarCadastro(w http.ResponseWriter, _ *http.Request) {
	err := templates.ExecuteTemplate(w, "cadastro.html", "a")
	if err != nil {
		return
	}
}

func autenticaCadastroELevaAoLogin(w http.ResponseWriter, r *http.Request) {
	nomecompleto := r.FormValue("nome_completo")
	cpf := r.FormValue("cpf")
	cns := r.FormValue("cns")
	cbo := r.FormValue("cbo")
	cnes := r.FormValue("cnes")
	ine := r.FormValue("ine")
	senha := r.FormValue("senha")
	confirmsenha := r.FormValue("confirmsenha")

	if confirmsenha == senha {
		_, err := db.Exec("INSERT INTO cadastro(nome_completo, cpf, cns, cbo, cnes, ine, senha) VALUES($1, $2, $3, $4, $5, $6, $7)", nomecompleto, cpf, cns, cbo, cnes, ine, senha)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	} else {
		http.Redirect(w, r, "/cadastro", http.StatusSeeOther)
		return
	}
	err := templates.ExecuteTemplate(w, "login.html", loginInvalido)
	if err != nil {
		return
	}
}

func colocarDados() {
	pegardados, err := db.Query("SELECT homem, etilista, tabagista, lesao_bucal, ultima_visita FROM pacientes")
	if err != nil {
		return
	}
	defer pegardados.Close()
	armazenamento := make([]PegarDados, 0)

	for pegardados.Next() {
		armazenar := PegarDados{}
		err := pegardados.Scan(&armazenar.Homem, &armazenar.Etilista, &armazenar.Tabagista, &armazenar.LesaoBucal, &armazenar.UltimaVisita)
		if err != nil {
			log.Println(err)
			return
		}
		armazenamento = append(armazenamento, armazenar)
	}
	if err = pegardados.Err(); err != nil {
		return
	}
	pgbaixo := &qtdBaixo
	pgmedio := &qtdMedio
	pgalto := &qtdAlto
	pgtotal := &qtdTotal
	for _, armazenado := range armazenamento {
		if armazenado.Tabagista == "Não" && armazenado.LesaoBucal == "Não" {
			*pgbaixo++
		} else if armazenado.Tabagista == "Sim" && armazenado.LesaoBucal == "Não" {
			*pgmedio++
		} else if armazenado.LesaoBucal == "Sim" {
			*pgalto++
		}
		quebrarUltimaVisita := strings.Split(armazenado.UltimaVisita, "-")
		armazenarQtdUltimaVisita := &qtdVisitadosMaisDeUmMes
		now := time.Now()
		diaVisita, _ := strconv.Atoi(quebrarUltimaVisita[2])
		mesVisita, _ := strconv.Atoi(quebrarUltimaVisita[1])
		if diaVisita <= now.Day() && mesVisita < int(now.Month()){
			*armazenarQtdUltimaVisita++
		} else if (int(now.Month()) - mesVisita) >= 2{
			*armazenarQtdUltimaVisita++
		} 
	}
	*pgtotal = *pgbaixo + *pgmedio + *pgalto
}

func loginInvalidado(w http.ResponseWriter, r *http.Request){
	err := templates.ExecuteTemplate(w, "login.html", loginInvalido)
	if err != nil{
		return
	}
	ponteiroLoginInvalido := &loginInvalido
	*ponteiroLoginInvalido = false
}

func autenticaLoginELevaAoDashboard(w http.ResponseWriter, r *http.Request) {
	ponteiroLoginInvalido := &loginInvalido
	*ponteiroLoginInvalido = false
	var endereco string
	cpf := &cpfLogin
	senha := &senhaLogin
	*cpf = r.FormValue("cpf")
	*senha = r.FormValue("senha")
	cpfsenha, err := db.Query("SELECT nome_completo, cpf, senha, cns, cbo, cnes, ine FROM cadastro")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	cepEndereco, err := db.Query("SELECT nome_completo, data_nasc, telefone, homem, etilista, tabagista, lesao_bucal, ultima_visita, cep, bairro, rua, numero, complemento FROM pacientes")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	defer cepEndereco.Close()
	defer cpfsenha.Close()
	armazenamento := make([]validarlogin, 0)

	for cpfsenha.Next() {
		armazenar := validarlogin{}
		err := cpfsenha.Scan(&armazenar.Usuario, &armazenar.Cpf, &armazenar.Senha, &armazenar.Cns, &armazenar.Cbo, &armazenar.Cnes, &armazenar.Ine)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		armazenar.QtdMaisDeUmMes = qtdVisitadosMaisDeUmMes
		armazenar.QtdBaixo = qtdBaixo
		armazenar.QtdMedio = qtdMedio
		armazenar.QtdAlto = qtdAlto
		armazenar.QtdTotal = qtdBaixo + qtdMedio + qtdAlto + 3
		if qtdBaixo == 0 && qtdMedio == 0 && qtdAlto == 0 {
			var porcbaixo float64 = 0
			var porcmedio float64 = 0
			var porcalto float64 = 0
			armazenar.PorcBaixo = porcbaixo
			armazenar.PorcMedio = porcmedio
			armazenar.PorcAlto = porcalto
		} else {
			porcbaixo := (float64(qtdBaixo) / float64(qtdTotal)) * 100
			porcmedio := (float64(qtdMedio) / float64(qtdTotal)) * 100
			porcalto := (float64(qtdAlto) / float64(qtdTotal)) * 100
			porcbaixo = float64(int(porcbaixo*100)) / 100
			porcmedio = float64(int(porcmedio*100)) / 100
			porcalto = float64(int(porcalto*100)) / 100
			armazenar.PorcBaixo = porcbaixo
			armazenar.PorcMedio = porcmedio
			armazenar.PorcAlto = porcalto
		}
		armazenamento = append(armazenamento, armazenar)
	}
	if err = cpfsenha.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	for _, armazenado := range armazenamento {
		if armazenado.Cpf == cpfLogin && armazenado.Senha == senhaLogin {
			for cepEndereco.Next() {
				armazenar := validarlogin{}
				armazenarDadosMaps := DadosGoogleMaps{}
				err = cepEndereco.Scan(&armazenarDadosMaps.Nome, &armazenarDadosMaps.DataNascimento, &armazenarDadosMaps.Telefone, &armazenar.Homem, &armazenar.Etilista, &armazenar.Tabagista, &armazenar.FeridasBucais, &armazenarDadosMaps.UltimaVisita, &armazenarDadosMaps.Cep, &armazenar.Bairro, &armazenar.Rua, &armazenar.Numero, &armazenar.Complemento)
				quebrar := strings.Split(armazenarDadosMaps.UltimaVisita, "-")
				ultimaVisita := quebrar[2] + "/" + quebrar[1] + "/" + quebrar[0]
				diaUltimaVisita, _ := strconv.Atoi(quebrar[2])
				mesUltimaVisita, _ := strconv.Atoi(quebrar[1])
				armazenarDadosMaps.UltimaVisita = ultimaVisita
				now := time.Now()
				if mesUltimaVisita < int(now.Month()) && diaUltimaVisita <= now.Day(){
					armazenarDadosMaps.MaisDeUmMes = true
				} else if (int(now.Month()) - mesUltimaVisita) >= 2{
					armazenarDadosMaps.MaisDeUmMes = true
				} else{
					armazenarDadosMaps.MaisDeUmMes = false
				}
				if armazenar.Complemento != "" {
					endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Complemento + "," + armazenar.Bairro
				} else {
					endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro
				}
				armazenarDadosMaps.Endereco = endereco
				if armazenar.Tabagista == "Não" && armazenar.FeridasBucais == "Não" {
					armazenarDadosMaps.Baixo = true
					armazenarDadosMaps.Medio = false
					armazenarDadosMaps.Alto = false
					if armazenar.Etilista == "Sim" && armazenar.Homem == "Sim" {
						fatores := "Homem, Etilista"
						armazenarDadosMaps.Fatores = fatores
					} else if armazenar.Etilista == "Sim" && armazenar.Homem == "Não" {
						fatores := "Etilista"
						armazenarDadosMaps.Fatores = fatores
					} else if armazenar.Etilista == "Não" && armazenar.Homem == "Sim" {
						fatores := "Homem"
						armazenarDadosMaps.Fatores = fatores
					}
				} else if armazenar.Tabagista == "Sim" && armazenar.FeridasBucais == "Não" {
					armazenarDadosMaps.Baixo = false
					armazenarDadosMaps.Medio = true
					armazenarDadosMaps.Alto = false
					if armazenar.Etilista == "Sim" && armazenar.Homem == "Sim" {
						fatores := "Homem, Etilista, Tabagista"
						armazenarDadosMaps.Fatores = fatores
					} else if armazenar.Etilista == "Sim" && armazenar.Homem == "Não" {
						fatores := "Etilista, Tabagista"
						armazenarDadosMaps.Fatores = fatores
					} else if armazenar.Etilista == "Não" && armazenar.Homem == "Sim" {
						fatores := "Homem, Tabagista"
						armazenarDadosMaps.Fatores = fatores
					}
				} else if armazenar.FeridasBucais == "Sim" {
					armazenarDadosMaps.Baixo = false
					armazenarDadosMaps.Medio = false
					armazenarDadosMaps.Alto = true
					if armazenar.Etilista == "Sim" && armazenar.Homem == "Sim" && armazenar.Tabagista == "Sim" {
						fatores := "Homem, Etilista, Tabagista, Feridas Bucais"
						armazenarDadosMaps.Fatores = fatores
					} else if armazenar.Etilista == "Sim" && armazenar.Homem == "Não" && armazenar.Tabagista == "Sim" {
						fatores := "Etilista, Tabagista, Feridas Bucais"
						armazenarDadosMaps.Fatores = fatores
					} else if armazenar.Etilista == "Não" && armazenar.Homem == "Sim" && armazenar.Tabagista == "Sim" {
						fatores := "Homem, Tabagista, Feridas Bucais"
						armazenarDadosMaps.Fatores = fatores
					} else if armazenar.Etilista == "Sim" && armazenar.Homem == "Sim" && armazenar.Tabagista == "Não" {
						fatores := "Homem, Etilista, Feridas Bucais"
						armazenarDadosMaps.Fatores = fatores
					} else if armazenar.Homem == "Não" && armazenar.Etilista == "Sim" && armazenar.Tabagista == "Não" {
						fatores := "Etilista, Feridas Bucais"
						armazenarDadosMaps.Fatores = fatores
					} else if armazenar.Homem == "Não" && armazenar.Etilista == "Não" && armazenar.Tabagista == "Não" {
						fatores := "Feridas Bucais"
						armazenarDadosMaps.Fatores = fatores
					}
				}
				armazenado.DadosGoogleMaps = append(armazenado.DadosGoogleMaps, armazenarDadosMaps)
				if err != nil {
					log.Println(err)
					http.Error(w, http.StatusText(500), 500)
					return
				}
			}
			armazenador := &primeiraletraLogin
			armazenador2 := &usuarioLogin
			armazenado.PrimeiraLetra = string(armazenado.Usuario[0])
			*armazenador = string(armazenado.Usuario[0])
			quebrado := strings.Split(armazenado.Usuario, " ")
			armazenado.Usuario = quebrado[0]
			*armazenador2 = armazenado.Usuario
			cns := &Cns
			cbo := &Cbo
			cnes := &Cnes
			ine := &Ine
			*ine = armazenado.Ine
			*cns = armazenado.Cns
			*cnes = armazenado.Cnes
			*cbo = armazenado.Cbo
			err = templates.ExecuteTemplate(w, "dashboard.html", armazenado)
			if err != nil {
				return
			}
			return
		}
	}
	ponteiroLoginInvalido = &loginInvalido
	*ponteiroLoginInvalido = true
	http.Redirect(w, r, "/login/autenticar", http.StatusSeeOther)
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	armazenado := validarlogin{}
	var endereco string
	var porcbaixo, porcmedio, porcalto float64
	if qtdBaixo == 0 && qtdMedio == 0 && qtdAlto == 0 {
		porcbaixo = 0
		porcmedio = 0
		porcalto = 0
		armazenado.PorcBaixo = porcbaixo
		armazenado.PorcMedio = porcmedio
		armazenado.PorcAlto = porcalto
	} else {
		porcbaixo = (float64(qtdBaixo) / float64(qtdTotal)) * 100
		porcmedio = (float64(qtdMedio) / float64(qtdTotal)) * 100
		porcalto = (float64(qtdAlto) / float64(qtdTotal)) * 100
		porcbaixo = float64(int(porcbaixo*100)) / 100
		porcmedio = float64(int(porcmedio*100)) / 100
		porcalto = float64(int(porcalto*100)) / 100
		armazenado.PorcBaixo = porcbaixo
		armazenado.PorcMedio = porcmedio
		armazenado.PorcAlto = porcalto
	}
	armazenado = validarlogin{Usuario: usuarioLogin, PrimeiraLetra: primeiraletraLogin, QtdBaixo: qtdBaixo, QtdMedio: qtdMedio, QtdAlto: qtdAlto, QtdMaisDeUmMes: qtdVisitadosMaisDeUmMes, PorcBaixo: porcbaixo, PorcMedio: porcmedio, PorcAlto: porcalto}
	ponteiroConfirmCadastro := &confirmCadastro
	*ponteiroConfirmCadastro = false
	ponteiroErroCampos := &erroCadastro
	*ponteiroErroCampos = false
	cepEndereco, err := db.Query("SELECT nome_completo, data_nasc, telefone, homem, etilista, tabagista, lesao_bucal, ultima_visita, cep, bairro, rua, numero, complemento FROM pacientes")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	defer cepEndereco.Close()
	for cepEndereco.Next() {
		armazenar := validarlogin{}
		armazenarDadosMaps := DadosGoogleMaps{}
		err = cepEndereco.Scan(&armazenarDadosMaps.Nome, &armazenarDadosMaps.DataNascimento, &armazenarDadosMaps.Telefone, &armazenar.Homem, &armazenar.Etilista, &armazenar.Tabagista, &armazenar.FeridasBucais, &armazenarDadosMaps.UltimaVisita, &armazenarDadosMaps.Cep, &armazenar.Bairro, &armazenar.Rua, &armazenar.Numero, &armazenar.Complemento)
		quebrar := strings.Split(armazenarDadosMaps.UltimaVisita, "-")
		ultimaVisita := quebrar[2] + "/" + quebrar[1] + "/" + quebrar[0]
		armazenarDadosMaps.UltimaVisita = ultimaVisita
		diaUltimaVisita, _ := strconv.Atoi(quebrar[2])
		mesUltimaVisita, _ := strconv.Atoi(quebrar[1])
		now := time.Now()
		if mesUltimaVisita < int(now.Month()) && diaUltimaVisita <= now.Day(){
			armazenarDadosMaps.MaisDeUmMes = true
		} else if (int(now.Month()) - mesUltimaVisita) >= 2{
			armazenarDadosMaps.MaisDeUmMes = true
		} else{
			armazenarDadosMaps.MaisDeUmMes = false
		}
		if armazenar.Complemento != "" {
			endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Complemento + "," + armazenar.Bairro
		} else {
			endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro
		}
		armazenarDadosMaps.Endereco = endereco
		if armazenar.Tabagista == "Não" && armazenar.FeridasBucais == "Não" {
			armazenarDadosMaps.Baixo = true
			armazenarDadosMaps.Medio = false
			armazenarDadosMaps.Alto = false
			if armazenar.Etilista == "Sim" && armazenar.Homem == "Sim" {
				fatores := "Homem, Etilista"
				armazenarDadosMaps.Fatores = fatores
			} else if armazenar.Etilista == "Sim" && armazenar.Homem == "Não" {
				fatores := "Etilista"
				armazenarDadosMaps.Fatores = fatores
			} else if armazenar.Etilista == "Não" && armazenar.Homem == "Sim" {
				fatores := "Homem"
				armazenarDadosMaps.Fatores = fatores
			}
		} else if armazenar.Tabagista == "Sim" && armazenar.FeridasBucais == "Não" {
			armazenarDadosMaps.Baixo = false
			armazenarDadosMaps.Medio = true
			armazenarDadosMaps.Alto = false
			if armazenar.Etilista == "Sim" && armazenar.Homem == "Sim" {
				fatores := "Homem, Etilista, Tabagista"
				armazenarDadosMaps.Fatores = fatores
			} else if armazenar.Etilista == "Sim" && armazenar.Homem == "Não" {
				fatores := "Etilista, Tabagista"
				armazenarDadosMaps.Fatores = fatores
			} else if armazenar.Etilista == "Não" && armazenar.Homem == "Sim" {
				fatores := "Homem, Tabagista"
				armazenarDadosMaps.Fatores = fatores
			}
		} else if armazenar.FeridasBucais == "Sim" {
			armazenarDadosMaps.Baixo = false
			armazenarDadosMaps.Medio = false
			armazenarDadosMaps.Alto = true
			if armazenar.Etilista == "Sim" && armazenar.Homem == "Sim" && armazenar.Tabagista == "Sim" {
				fatores := "Homem, Etilista, Tabagista, Feridas Bucais"
				armazenarDadosMaps.Fatores = fatores
			} else if armazenar.Etilista == "Sim" && armazenar.Homem == "Não" && armazenar.Tabagista == "Sim" {
				fatores := "Etilista, Tabagista, Feridas Bucais"
				armazenarDadosMaps.Fatores = fatores
			} else if armazenar.Etilista == "Não" && armazenar.Homem == "Sim" && armazenar.Tabagista == "Sim" {
				fatores := "Homem, Tabagista, Feridas Bucais"
				armazenarDadosMaps.Fatores = fatores
			} else if armazenar.Etilista == "Sim" && armazenar.Homem == "Sim" && armazenar.Tabagista == "Não" {
				fatores := "Homem, Etilista, Feridas Bucais"
				armazenarDadosMaps.Fatores = fatores
			} else if armazenar.Homem == "Não" && armazenar.Etilista == "Sim" && armazenar.Tabagista == "Não" {
				fatores := "Etilista, Feridas Bucais"
				armazenarDadosMaps.Fatores = fatores
			} else if armazenar.Homem == "Não" && armazenar.Etilista == "Não" && armazenar.Tabagista == "Não" {
				fatores := "Feridas Bucais"
				armazenarDadosMaps.Fatores = fatores
			}
		}
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		armazenado.DadosGoogleMaps = append(armazenado.DadosGoogleMaps, armazenarDadosMaps)
	}
	armazenado.QtdTotal = qtdAlto + qtdBaixo + qtdMedio + 3
	err = templates.ExecuteTemplate(w, "dashboard.html", armazenado)
	if err != nil {
		return
	}
}

func executarEsqueceuSenha(w http.ResponseWriter, _ *http.Request) {
	err := templates.ExecuteTemplate(w, "esqueceusenha.html", esqueceuInvalido)
	ponteiroEsqueceuInvalido := &esqueceuInvalido
	*ponteiroEsqueceuInvalido = false
	if err != nil {
		return
	}
}

func atualizarSenha(w http.ResponseWriter, r *http.Request) {
	ponteiroEsqueceuInvalido := &esqueceuInvalido
	*ponteiroEsqueceuInvalido = false
	cpf := r.FormValue("cpf")
	senha := r.FormValue("senha")
	confirmarsenha := r.FormValue("confirmpassword")
	pegarcpf, err := db.Query("SELECT cpf FROM cadastro")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	defer pegarcpf.Close()
	armazenamento := make([]validarCpf, 0)

	for pegarcpf.Next() {
		armazenar := validarCpf{}
		err := pegarcpf.Scan(&armazenar.Cpf)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		armazenamento = append(armazenamento, armazenar)
	}
	if err = pegarcpf.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	for _, armazenado := range armazenamento {
		if armazenado.Cpf == cpf && senha == confirmarsenha {
			_, err := db.Exec(`UPDATE cadastro SET senha=$1 WHERE cpf=$2`, senha, cpf)
			if err != nil {
				return
			}
			err = templates.ExecuteTemplate(w, "login.html", loginInvalido)
			ponteiroLoginInvalido := &loginInvalido
			*ponteiroLoginInvalido = false
			if err != nil {
				return
			}
			return
		}
	}
	ponteiroEsqueceuInvalido = &esqueceuInvalido
	*ponteiroEsqueceuInvalido = true
	http.Redirect(w, r, "/esqueceusenha/invalido", http.StatusSeeOther)
}

func esqueceuSenhaInvalidado(w http.ResponseWriter, r *http.Request){
	err := templates.ExecuteTemplate(w, "esqueceusenha.html", esqueceuInvalido)
	if err != nil{
		return
	}
	ponteiroEsqueceuInvalido := &esqueceuInvalido
	*ponteiroEsqueceuInvalido = false
}

func executarCentralUsuario(w http.ResponseWriter, r *http.Request) {
	cpfsenha, err := db.Query("SELECT nome_completo, cpf, cns, cbo, cnes, ine, senha FROM cadastro")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	defer cpfsenha.Close()
	armazenamento := make([]ACS, 0)

	for cpfsenha.Next() {
		armazenar := ACS{}
		err := cpfsenha.Scan(&armazenar.NomeCompleto, &armazenar.CPF, &armazenar.CNS, &armazenar.CBO, &armazenar.CNES, &armazenar.INE, &armazenar.SenhaACS)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		armazenamento = append(armazenamento, armazenar)
	}
	if err = cpfsenha.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	for _, armazenado := range armazenamento {
		if armazenado.CPF == cpfLogin && armazenado.SenhaACS == senhaLogin {
			armazenado.PrimeiraLetra = string(armazenado.NomeCompleto[0])
			armazenado.CPF = strings.ReplaceAll(armazenado.CPF, armazenado.CPF[:5], "*****")
			armazenado.CNS = strings.ReplaceAll(armazenado.CNS, armazenado.CNS[:5], "*****")
			armazenado.CNES = strings.ReplaceAll(armazenado.CNES, armazenado.CNES[:3], "***")
			quebrado2 := strings.Split(armazenado.NomeCompleto, " ")
			armazenado.User = quebrado2[0]
			quebrado := strings.Split(armazenado.SenhaACS, "")
			for i := 0; i < len(quebrado); i++ {
				armazenado.SenhaACS = strings.Replace(armazenado.SenhaACS, quebrado[i], "*", -1)
			}
			err = templates.ExecuteTemplate(w, "centralusuario.html", armazenado)
			if err != nil {
				return
			}
			return
		}
	}

}

func atualizarSenhaCentralUsuario(w http.ResponseWriter, r *http.Request) {
	novasenha := r.FormValue("senha")
	_, err := db.Exec(`UPDATE cadastro SET senha=$1 WHERE cpf=$2`, novasenha, cpfLogin)
	if err != nil {
		return
	}
	senhaLogin = novasenha
	cpfsenha, err := db.Query("SELECT nome_completo, cpf, cns, cbo, cnes, ine, senha FROM cadastro")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	defer cpfsenha.Close()
	armazenamento := make([]ACS, 0)

	for cpfsenha.Next() {
		armazenar := ACS{}
		err := cpfsenha.Scan(&armazenar.NomeCompleto, &armazenar.CPF, &armazenar.CNS, &armazenar.CBO, &armazenar.CNES, &armazenar.INE, &armazenar.SenhaACS)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		armazenamento = append(armazenamento, armazenar)
	}
	if err = cpfsenha.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	for _, armazenado := range armazenamento {
		if armazenado.CPF == cpfLogin && armazenado.SenhaACS == senhaLogin {
			armazenado.PrimeiraLetra = string(armazenado.NomeCompleto[0])
			armazenado.CPF = strings.ReplaceAll(armazenado.CPF, armazenado.CPF[:5], "*****")
			armazenado.CNS = strings.ReplaceAll(armazenado.CNS, armazenado.CNS[:5], "*****")
			armazenado.CNES = strings.ReplaceAll(armazenado.CNES, armazenado.CNES[:3], "***")
			quebrado2 := strings.Split(armazenado.NomeCompleto, " ")
			armazenado.User = quebrado2[0]
			quebrado := strings.Split(armazenado.SenhaACS, "")
			for i := 0; i < len(quebrado); i++ {
				armazenado.SenhaACS = strings.Replace(armazenado.SenhaACS, quebrado[i], "*", -1)
			}
			err = templates.ExecuteTemplate(w, "centralusuario.html", armazenado)
			if err != nil {
				return
			}
			return
		}
	}
}

func executarFormulario(w http.ResponseWriter, _ *http.Request) {
	cnsq := strings.Split(Cns, "")
	cboq := strings.Split(Cbo, "")
	cnesq := strings.Split(Cnes, "")
	ineq := strings.Split(Ine, "")
	cbo1 := cboq[0]
	cbo2 := cboq[1]
	cbo3 := cboq[2]
	cbo4 := cboq[3]
	cbo5 := cboq[4]
	cbo6 := cboq[5]
	d := DadosForm{Usuario: usuarioLogin, PrimeiraLetra: primeiraletraLogin, Cns: cnsq, Cbo1: cbo1, Cbo2: cbo2, Cbo3: cbo3, Cbo4: cbo4, Cbo5: cbo5, Cbo6: cbo6, Cnes: cnesq, Ine: ineq, Confirmacao: confirmCadastro, ErroCampos: erroCadastro}
	err := templates.ExecuteTemplate(w, "formulario.html", d)
	if err != nil {
		return
	}
}

func cadastrarPaciente(w http.ResponseWriter, r *http.Request) {
	nome := r.FormValue("nome")
	datanascimento := r.FormValue("datanascimento")
	cpf := r.FormValue("cpfpaciente")
	nomemae := r.FormValue("nomemae")
	sexo := r.FormValue("sexo")
	cartaosus := r.FormValue("cartaosus")
	telefone := r.FormValue("telefone")
	email := r.FormValue("email")
	cep := r.FormValue("cep")
	bairro := r.FormValue("bairro")
	rua := r.FormValue("rua")
	numero, _ := strconv.Atoi(r.FormValue("numero"))
	complemento := r.FormValue("complemento")
	homem := r.FormValue("tipo1")
	etilista := r.FormValue("tipo2")
	tabagista := r.FormValue("tipo3")
	lesao_bucal := r.FormValue("tipo4")
	data_cadastro := r.FormValue("DataCadastro")
	cnsq := strings.Split(Cns, "")
	cboq := strings.Split(Cbo, "")
	cnesq := strings.Split(Cnes, "")
	ineq := strings.Split(Ine, "")
	cbo1 := cboq[0]
	cbo2 := cboq[1]
	cbo3 := cboq[2]
	cbo4 := cboq[3]
	cbo5 := cboq[4]
	cbo6 := cboq[5]
	d := DadosForm{Usuario: usuarioLogin, PrimeiraLetra: primeiraletraLogin, Cns: cnsq, Cbo1: cbo1, Cbo2: cbo2, Cbo3: cbo3, Cbo4: cbo4, Cbo5: cbo5, Cbo6: cbo6, Cnes: cnesq, Ine: ineq}

	if homem != "" && etilista != "" && tabagista != "" && lesao_bucal != "" && sexo != "" && (homem != "Não" || etilista != "Não" || tabagista != "Não" || lesao_bucal != "Não") {
		_, err := db.Exec("INSERT INTO pacientes(nome_completo, data_nasc, cpf, nome_mae, sexo, cartao_sus, telefone, email, cep, bairro, rua, numero, complemento, homem, etilista, tabagista, lesao_bucal, data_cadastro, ultima_visita) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)", nome, datanascimento, cpf, nomemae, sexo, cartaosus, telefone, email, cep, bairro, rua, numero, complemento, homem, etilista, tabagista, lesao_bucal, data_cadastro, data_cadastro)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		ponteiroConfirmando := &confirmCadastro
		quebrarDataCadastro := strings.Split(data_cadastro, "-")
		now := time.Now()
		armazenarQtdVisitadosMaisDeUmMes := &qtdVisitadosMaisDeUmMes
		diaCadastro, _ := strconv.Atoi(quebrarDataCadastro[2])
		mesCadastro, _ := strconv.Atoi(quebrarDataCadastro[1])
		if mesCadastro < int(now.Month()) && diaCadastro <= now.Day(){
			*armazenarQtdVisitadosMaisDeUmMes++
		} else if (int(now.Month()) - mesCadastro) >= 2{
			*armazenarQtdVisitadosMaisDeUmMes++
		}
		*ponteiroConfirmando = true
		ponteiroErro := &erroCadastro
		*ponteiroErro = false
		d.Confirmacao = confirmCadastro
		d.ErroCampos = erroCadastro
		err = templates.ExecuteTemplate(w, "formulario.html", d)
		if err != nil {
			return
		}
		pgbaixo := &qtdBaixo
		pgmedio := &qtdMedio
		pgalto := &qtdAlto
		pgtotal := &qtdTotal
		if tabagista == "Não" && lesao_bucal == "Não" {
			*pgbaixo++
		} else if tabagista == "Sim" && lesao_bucal == "Não" {
			*pgmedio++
		} else if lesao_bucal == "Sim" {
			*pgalto++
		}
		*pgtotal = *pgalto + *pgmedio + *pgbaixo
	} else {
		ponteiroConfirmando := &confirmCadastro
		*ponteiroConfirmando = false
		ponteiroErro := &erroCadastro
		*ponteiroErro = true
		http.Redirect(w, r, "/cadastrar-paciente", http.StatusSeeOther)
	}
}

func executarPagFaq(w http.ResponseWriter, _ *http.Request) {
	u := UsuarioNoDashboard{Usuario: usuarioLogin, Primeira: primeiraletraLogin}
	err := templates.ExecuteTemplate(w, "faq.html", u)
	if err != nil {
		return
	}
}

func executarPacGerais(w http.ResponseWriter, _ *http.Request) {
	u := UsuarioNoDashboard{Usuario: usuarioLogin, Primeira: primeiraletraLogin, QtdBaixo: qtdBaixo, QtdMedio: qtdMedio, QtdAlto: qtdAlto}
	err := templates.ExecuteTemplate(w, "pacGerais.html", u)
	if err != nil {
		return
	}
}

func executarPgBaixo(w http.ResponseWriter, _ *http.Request) {
	pesquisa, err := db.Query("SELECT nome_completo, data_nasc, telefone, bairro, rua, numero, complemento, homem, etilista, tabagista, lesao_bucal FROM pacientes ORDER BY nome_completo")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	defer pesquisa.Close()
	var armazenamento []Pacientes
	var temDados bool
	for pesquisa.Next() {
		armazenar := Pacientes{}
		err := pesquisa.Scan(&armazenar.Nome, &armazenar.DataNasc, &armazenar.Telefone, &armazenar.Bairro, &armazenar.Rua, &armazenar.Numero, &armazenar.Complemento, &armazenar.Homem, &armazenar.Etilista, &armazenar.Tabagista, &armazenar.LesaoBucal)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		quebrar := strings.Split(armazenar.DataNasc, "/")
		if armazenar.Complemento != "" {
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro + "," + armazenar.Complemento
		} else {
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro
		}
		if armazenar.Tabagista == "Não" && armazenar.LesaoBucal == "Não" {
			if armazenar.Homem == "Sim" && armazenar.Etilista == "Sim" {
				armazenar.Fatores = "Homem/Etilista"
			} else if armazenar.Homem == "Sim" && armazenar.Etilista == "Não" {
				armazenar.Fatores = "Homem"
			} else if armazenar.Homem == "Não" && armazenar.Etilista == "Sim" {
				armazenar.Fatores = "Mulher/Etilista"
			}
			now := time.Now()
			dia, _ := strconv.Atoi(quebrar[0])
			mes, _ := strconv.Atoi(quebrar[1])
			ano, _ := strconv.Atoi(quebrar[2])
			armazenar.Idade = now.Year() - ano
			if int(now.Month()) < mes || (int(now.Month()) == mes && now.Day() < dia) {
				armazenar.Idade--
			}
			armazenar.TemDados = true
			armazenar.Usuario = usuarioLogin
			armazenar.Primeira = primeiraletraLogin
			armazenamento = append(armazenamento, armazenar)
			temDados = true
		}
	}
	if !temDados {
		u := Pacientes{}
		u.Usuario = usuarioLogin
		u.Primeira = primeiraletraLogin
		u.TemDados = false
		armazenamento = append(armazenamento, u)
	}
	err = templates.ExecuteTemplate(w, "pg-baixo.html", armazenamento)
	if err != nil {
		return
	}
}

func executarPgBaixoFiltro(w http.ResponseWriter, r *http.Request) {
	masculino := r.FormValue("radio1")
	feminino := r.FormValue("radio2")
	idade1 := r.FormValue("radio3")
	idade2 := r.FormValue("radio4")
	idade3 := r.FormValue("radio5")
	idade4 := r.FormValue("radio6")
	etilista := r.FormValue("radio7")
	pesquisa, err := db.Query("SELECT nome_completo, data_nasc, sexo, telefone, bairro, rua, numero, complemento, homem, etilista, tabagista, lesao_bucal FROM pacientes ORDER BY nome_completo")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	defer pesquisa.Close()
	var armazenamento []Pacientes
	for pesquisa.Next() {
		armazenar := Pacientes{}
		err := pesquisa.Scan(&armazenar.Nome, &armazenar.DataNasc, &armazenar.Sexo, &armazenar.Telefone, &armazenar.Bairro, &armazenar.Rua, &armazenar.Numero, &armazenar.Complemento, &armazenar.Homem, &armazenar.Etilista, &armazenar.Tabagista, &armazenar.LesaoBucal)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		quebrar := strings.Split(armazenar.DataNasc, "/")
		if armazenar.Complemento != "" {
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro + "," + armazenar.Complemento
		} else {
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro
		}
		if armazenar.Tabagista == "Não" && armazenar.LesaoBucal == "Não" {
			if armazenar.Homem == "Sim" && armazenar.Etilista == "Sim" {
				armazenar.Fatores = "Homem/Etilista"
			} else if armazenar.Homem == "Sim" && armazenar.Etilista == "Não" {
				armazenar.Fatores = "Homem"
			} else if armazenar.Homem == "Não" && armazenar.Etilista == "Sim" {
				armazenar.Fatores = "Mulher/Etilista"
			}
			now := time.Now()
			dia, _ := strconv.Atoi(quebrar[0])
			mes, _ := strconv.Atoi(quebrar[1])
			ano, _ := strconv.Atoi(quebrar[2])
			armazenar.Idade = now.Year() - ano
			if int(now.Month()) < mes || (int(now.Month()) == mes && now.Day() < dia) {
				armazenar.Idade--
			}
			armazenar.Usuario = usuarioLogin
			armazenar.Primeira = primeiraletraLogin
			armazenamento = append(armazenamento, armazenar)
		}
	}
	var armazenadoPgBaixo []Pacientes
	armazenar2 := Pacientes{}
	armazenar2.Usuario = usuarioLogin
	armazenar2.Primeira = primeiraletraLogin
	armazenadoPgBaixo = append(armazenadoPgBaixo, armazenar2)
	for _, armazenado := range armazenamento {
		if masculino == "Masculino" && armazenado.Sexo == "Masculino" {
			armazenado.TemDados = true
			if etilista == "Etilista" && armazenado.Etilista == "Sim" {
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
			}
			if etilista == "" && armazenado.Etilista == "Não" {
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
			}
		}
		if feminino == "Feminino" && armazenado.Sexo == "Feminino" {
			armazenado.TemDados = true
			if etilista == "Etilista" && armazenado.Etilista == "Sim" {
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
			}
			if etilista == "" && armazenado.Etilista == "Não" {
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
			}
		}
		if feminino == "" && masculino == "" {
			if etilista == "Etilista" && armazenado.Etilista == "Sim" {
				armazenado.TemDados = true
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
			} else if etilista == "" && armazenado.Etilista == "Sim" || armazenado.Homem == "Sim" {
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenado.TemDados = true
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenado.TemDados = true
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenado.TemDados = true
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenado.TemDados = true
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgBaixo = append(armazenadoPgBaixo, armazenado)
				}
			}
		}
	}
	err = templates.ExecuteTemplate(w, "pg-baixo.html", armazenadoPgBaixo)
	if err != nil {
		return
	}
}

func executarPgBaixoFiltroPorNome(w http.ResponseWriter, r *http.Request) {
	busca := strings.TrimSpace(r.FormValue("nome"))
	buscar, err := db.Query("SELECT nome_completo, data_nasc, sexo, telefone, bairro, rua, numero, complemento, homem, etilista, tabagista, lesao_bucal FROM pacientes WHERE nome_completo LIKE concat('%', text($1), '%')", busca)
	if err != nil {
		fmt.Println(err)
	}
	defer buscar.Close()
	var armazenamento []Pacientes
	armazenar2 := Pacientes{}
	armazenar2.Usuario = usuarioLogin
	armazenar2.Primeira = primeiraletraLogin
	armazenamento = append(armazenamento, armazenar2)
	for buscar.Next() {
		armazenar := Pacientes{}
		err = buscar.Scan(&armazenar.Nome, &armazenar.DataNasc, &armazenar.Sexo, &armazenar.Telefone, &armazenar.Bairro, &armazenar.Rua, &armazenar.Numero, &armazenar.Complemento, &armazenar.Homem, &armazenar.Etilista, &armazenar.Tabagista, &armazenar.LesaoBucal)
		if err != nil {
			panic(err.Error())
		}
		quebrar := strings.Split(armazenar.DataNasc, "/")
		if armazenar.Complemento != "" {
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro + "," + armazenar.Complemento
		} else {
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro
		}
		if armazenar.Tabagista == "Não" && armazenar.LesaoBucal == "Não" {
			if armazenar.Homem == "Sim" && armazenar.Etilista == "Sim" {
				armazenar.Fatores = "Homem/Etilista"
			} else if armazenar.Homem == "Sim" && armazenar.Etilista == "Não" {
				armazenar.Fatores = "Homem"
			} else if armazenar.Homem == "Não" && armazenar.Etilista == "Sim" {
				armazenar.Fatores = "Mulher/Etilista"
			}
			now := time.Now()
			dia, _ := strconv.Atoi(quebrar[0])
			mes, _ := strconv.Atoi(quebrar[1])
			ano, _ := strconv.Atoi(quebrar[2])
			armazenar.Idade = now.Year() - ano
			if int(now.Month()) < mes || (int(now.Month()) == mes && now.Day() < dia) {
				armazenar.Idade--
			}
			armazenar.TemDados = true
			armazenar.Usuario = usuarioLogin
			armazenar.Primeira = primeiraletraLogin
			armazenamento = append(armazenamento, armazenar)
		}
	}
	err = templates.ExecuteTemplate(w, "pg-baixo.html", armazenamento)
	if err != nil {
		return
	}
}

func executarPgMedio(w http.ResponseWriter, _ *http.Request) {
	pesquisa, err := db.Query("SELECT nome_completo, data_nasc, telefone, bairro, rua, numero, complemento, homem, etilista, tabagista, lesao_bucal FROM pacientes ORDER BY nome_completo")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	defer pesquisa.Close()
	var temDados bool
	var armazenamento []Pacientes
	for pesquisa.Next() {
		armazenar := Pacientes{}
		err := pesquisa.Scan(&armazenar.Nome, &armazenar.DataNasc, &armazenar.Telefone, &armazenar.Bairro, &armazenar.Rua, &armazenar.Numero, &armazenar.Complemento, &armazenar.Homem, &armazenar.Etilista, &armazenar.Tabagista, &armazenar.LesaoBucal)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		quebrar := strings.Split(armazenar.DataNasc, "/")
		if armazenar.Complemento != "" {
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro + "," + armazenar.Complemento
		} else {
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro
		}
		if armazenar.Tabagista == "Sim" && armazenar.LesaoBucal == "Não" {
			if armazenar.Homem == "Sim" && armazenar.Etilista == "Sim" && armazenar.Tabagista == "Sim" {
				armazenar.Fatores = "Homem/Etilista/Tabagista"
			} else if armazenar.Homem == "Sim" && armazenar.Etilista == "Não" && armazenar.Tabagista == "Sim" {
				armazenar.Fatores = "Homem/Tabagista"
			} else if armazenar.Homem == "Não" && armazenar.Tabagista == "Sim" && armazenar.Etilista == "Não" {
				armazenar.Fatores = "Mulher/Tabagista"
			} else if armazenar.Homem == "Não" && armazenar.Tabagista == "Sim" && armazenar.Etilista == "Sim" {
				armazenar.Fatores = "Mulher/Etilista/Tabagista"
			}
			now := time.Now()
			dia, _ := strconv.Atoi(quebrar[0])
			mes, _ := strconv.Atoi(quebrar[1])
			ano, _ := strconv.Atoi(quebrar[2])
			armazenar.Idade = now.Year() - ano
			if int(now.Month()) < mes || (int(now.Month()) == mes && now.Day() < dia) {
				armazenar.Idade--
			}
			armazenar.TemDados = true
			armazenar.Usuario = usuarioLogin
			armazenar.Primeira = primeiraletraLogin
			armazenamento = append(armazenamento, armazenar)
			temDados = true
		}
	}
	if !temDados {
		u := Pacientes{}
		u.Usuario = usuarioLogin
		u.Primeira = primeiraletraLogin
		u.TemDados = false
		armazenamento = append(armazenamento, u)
	}
	err = templates.ExecuteTemplate(w, "pg-medio.html", armazenamento)
	if err != nil {
		return
	}
}

func executarPgMedioFiltro(w http.ResponseWriter, r *http.Request) {
	masculino := r.FormValue("radio1")
	feminino := r.FormValue("radio2")
	idade1 := r.FormValue("radio3")
	idade2 := r.FormValue("radio4")
	idade3 := r.FormValue("radio5")
	idade4 := r.FormValue("radio6")
	etilista := r.FormValue("radio7")
	tabagista := r.FormValue("radio8")
	pesquisa, err := db.Query("SELECT nome_completo, data_nasc, sexo, telefone, bairro, rua, numero, complemento, homem, etilista, tabagista, lesao_bucal FROM pacientes ORDER BY nome_completo")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	defer pesquisa.Close()
	var armazenamento []Pacientes
	for pesquisa.Next() {
		armazenar := Pacientes{}
		err := pesquisa.Scan(&armazenar.Nome, &armazenar.DataNasc, &armazenar.Sexo, &armazenar.Telefone, &armazenar.Bairro, &armazenar.Rua, &armazenar.Numero, &armazenar.Complemento, &armazenar.Homem, &armazenar.Etilista, &armazenar.Tabagista, &armazenar.LesaoBucal)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		quebrar := strings.Split(armazenar.DataNasc, "/")
		if armazenar.Complemento != "" {
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro + "," + armazenar.Complemento
		} else {
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro
		}
		if armazenar.Tabagista == "Sim" && armazenar.LesaoBucal == "Não" {
			if armazenar.Homem == "Sim" && armazenar.Etilista == "Sim" && armazenar.Tabagista == "Sim" {
				armazenar.Fatores = "Homem/Etilista/Tabagista"
			} else if armazenar.Homem == "Sim" && armazenar.Etilista == "Não" && armazenar.Tabagista == "Sim" {
				armazenar.Fatores = "Homem/Tabagista"
			} else if armazenar.Homem == "Não" && armazenar.Tabagista == "Sim" && armazenar.Etilista == "Não" {
				armazenar.Fatores = "Mulher/Tabagista"
			} else if armazenar.Homem == "Não" && armazenar.Tabagista == "Sim" && armazenar.Etilista == "Sim" {
				armazenar.Fatores = "Mulher/Etilista/Tabagista"
			}
			now := time.Now()
			dia, _ := strconv.Atoi(quebrar[0])
			mes, _ := strconv.Atoi(quebrar[1])
			ano, _ := strconv.Atoi(quebrar[2])
			armazenar.Idade = now.Year() - ano
			if int(now.Month()) < mes || (int(now.Month()) == mes && now.Day() < dia) {
				armazenar.Idade--
			}
			armazenar.Usuario = usuarioLogin
			armazenar.Primeira = primeiraletraLogin
			armazenamento = append(armazenamento, armazenar)
		}
	}
	var armazenadoPgMedio []Pacientes
	armazenar2 := Pacientes{}
	armazenar2.Usuario = usuarioLogin
	armazenar2.Primeira = primeiraletraLogin
	armazenadoPgMedio = append(armazenadoPgMedio, armazenar2)
	for _, armazenado := range armazenamento {
		if masculino == "Masculino" && armazenado.Sexo == "Masculino" {
			armazenado.TemDados = true
			if tabagista == "Tabagista" && etilista == "Etilista" && armazenado.Tabagista == "Sim" && armazenado.Etilista == "Sim" {
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
			}
			if tabagista == "Tabagista" && etilista == "" && armazenado.Tabagista == "Sim" && armazenado.Etilista == "Não" {
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
			}
			if tabagista == "" && etilista == "" && armazenado.Tabagista == "Sim"{
				armazenado.TemDados = true
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
			}
		}
		if feminino == "Feminino" && armazenado.Sexo == "Feminino" {
			armazenado.TemDados = true
			if tabagista == "Tabagista" && etilista == "Etilista" && armazenado.Tabagista == "Sim" && armazenado.Etilista == "Sim" {
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
			}
			if tabagista == "Tabagista" && etilista == "" && armazenado.Tabagista == "Sim" && armazenado.Etilista == "Não" {
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
			}
			if tabagista == "" && etilista == "" && armazenado.Tabagista == "Sim"{
				armazenado.TemDados = true
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
			}
		}
		if feminino == "" && masculino == "" {
			if tabagista == "Tabagista" && etilista == "Etilista" && armazenado.Tabagista == "Sim" && armazenado.Etilista == "Sim" {
				armazenado.TemDados = true
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
			}
			if tabagista == "Tabagista" && etilista == "" && armazenado.Tabagista == "Sim" && armazenado.Etilista == "Não" {
				armazenado.TemDados = true
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
			}
			if tabagista == "" && etilista == "" && armazenado.Tabagista == "Sim"{
				armazenado.TemDados = true
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
			}
		}
	}
	err = templates.ExecuteTemplate(w, "pg-medio.html", armazenadoPgMedio)
	if err != nil {
		return
	}
}

func executarPgMedioFiltroPorNome(w http.ResponseWriter, r *http.Request) {
	busca := strings.TrimSpace(r.FormValue("nome"))
	buscar, err := db.Query("SELECT nome_completo, data_nasc, sexo, telefone, bairro, rua, numero, complemento, homem, etilista, tabagista, lesao_bucal FROM pacientes WHERE nome_completo LIKE concat('%', text($1), '%')", busca)
	if err != nil {
		fmt.Println(err)
	}
	defer buscar.Close()
	var armazenamento []Pacientes
	armazenar2 := Pacientes{}
	armazenar2.Usuario = usuarioLogin
	armazenar2.Primeira = primeiraletraLogin
	armazenamento = append(armazenamento, armazenar2)
	for buscar.Next() {
		armazenar := Pacientes{}
		err = buscar.Scan(&armazenar.Nome, &armazenar.DataNasc, &armazenar.Sexo, &armazenar.Telefone, &armazenar.Bairro, &armazenar.Rua, &armazenar.Numero, &armazenar.Complemento, &armazenar.Homem, &armazenar.Etilista, &armazenar.Tabagista, &armazenar.LesaoBucal)
		if err != nil {
			panic(err.Error())
		}
		quebrar := strings.Split(armazenar.DataNasc, "/")
		if armazenar.Complemento != "" {
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro + "," + armazenar.Complemento
		} else {
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro
		}
		if armazenar.Tabagista == "Sim" && armazenar.LesaoBucal == "Não" {
			if armazenar.Homem == "Sim" && armazenar.Etilista == "Sim" && armazenar.Tabagista == "Sim" {
				armazenar.Fatores = "Homem/Etilista/Tabagista"
			} else if armazenar.Homem == "Sim" && armazenar.Etilista == "Não" && armazenar.Tabagista == "Sim" {
				armazenar.Fatores = "Homem/Tabagista"
			} else if armazenar.Homem == "Não" && armazenar.Tabagista == "Sim" && armazenar.Etilista == "Não" {
				armazenar.Fatores = "Mulher/Tabagista"
			} else if armazenar.Homem == "Não" && armazenar.Tabagista == "Sim" && armazenar.Etilista == "Sim" {
				armazenar.Fatores = "Mulher/Etilista/Tabagista"
			}
			now := time.Now()
			dia, _ := strconv.Atoi(quebrar[0])
			mes, _ := strconv.Atoi(quebrar[1])
			ano, _ := strconv.Atoi(quebrar[2])
			armazenar.Idade = now.Year() - ano
			if int(now.Month()) < mes || (int(now.Month()) == mes && now.Day() < dia) {
				armazenar.Idade--
			}
			armazenar.TemDados = true
			armazenar.Usuario = usuarioLogin
			armazenar.Primeira = primeiraletraLogin
			armazenamento = append(armazenamento, armazenar)
		}
	}
	err = templates.ExecuteTemplate(w, "pg-medio.html", armazenamento)
	if err != nil {
		return
	}
}

func executarPgAlto(w http.ResponseWriter, _ *http.Request) {
	pesquisa, err := db.Query("SELECT nome_completo, data_nasc, sexo, telefone, bairro, rua, numero, complemento, homem, etilista, tabagista, lesao_bucal FROM pacientes ORDER BY nome_completo")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	defer pesquisa.Close()
	var armazenamento []Pacientes
	var temDados bool
	for pesquisa.Next() {
		armazenar := Pacientes{}
		err := pesquisa.Scan(&armazenar.Nome, &armazenar.DataNasc, &armazenar.Sexo, &armazenar.Telefone, &armazenar.Bairro, &armazenar.Rua, &armazenar.Numero, &armazenar.Complemento, &armazenar.Homem, &armazenar.Etilista, &armazenar.Tabagista, &armazenar.LesaoBucal)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		quebrar := strings.Split(armazenar.DataNasc, "/")
		if armazenar.Complemento != "" {
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro + "," + armazenar.Complemento
		} else {
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro
		}
		if armazenar.LesaoBucal == "Sim" {
			if armazenar.Homem == "Sim" {
				if armazenar.Etilista == "Sim" && armazenar.Tabagista == "Sim" {
					armazenar.Fatores = "Homem/Etilista/Tabagista/Feridas Bucais"
				} else if armazenar.Etilista == "Não" && armazenar.Tabagista == "Sim" {
					armazenar.Fatores = "Homem/Tabagista/Feridas Bucais"
				} else if armazenar.Etilista == "Sim" && armazenar.Tabagista == "Não" {
					armazenar.Fatores = "Homem/Etilista/Feridas Bucais"
				}
			}
			if armazenar.Homem == "Não" {
				if armazenar.Etilista == "Sim" && armazenar.Tabagista == "Sim" {
					armazenar.Fatores = "Mulher/Etilista/Tabagista/Feridas Bucais"
				} else if armazenar.Etilista == "Não" && armazenar.Tabagista == "Sim" {
					armazenar.Fatores = "Mulher/Tabagista/Feridas Bucais"
				} else if armazenar.Etilista == "Sim" && armazenar.Tabagista == "Não" {
					armazenar.Fatores = "Mulher/Etilista/Feridas Bucais"
				}
			}
			now := time.Now()
			dia, _ := strconv.Atoi(quebrar[0])
			mes, _ := strconv.Atoi(quebrar[1])
			ano, _ := strconv.Atoi(quebrar[2])
			armazenar.Idade = now.Year() - ano
			if int(now.Month()) < mes || (int(now.Month()) == mes && now.Day() < dia) {
				armazenar.Idade--
			}
			armazenar.TemDados = true
			armazenar.Usuario = usuarioLogin
			armazenar.Primeira = primeiraletraLogin
			armazenamento = append(armazenamento, armazenar)
			temDados = true
		}
	}
	if !temDados {
		u := Pacientes{}
		u.Usuario = usuarioLogin
		u.Primeira = primeiraletraLogin
		u.TemDados = false
		armazenamento = append(armazenamento, u)
	}
	err = templates.ExecuteTemplate(w, "pg-alto.html", armazenamento)
	if err != nil {
		return
	}
}

func executarPgAltoFiltro(w http.ResponseWriter, r *http.Request) {
	masculino := r.FormValue("radio1")
	feminino := r.FormValue("radio2")
	idade1 := r.FormValue("radio3")
	idade2 := r.FormValue("radio4")
	idade3 := r.FormValue("radio5")
	idade4 := r.FormValue("radio6")
	etilista := r.FormValue("radio7")
	tabagista := r.FormValue("radio8")
	feridasbucais := r.FormValue("radio9")
	pesquisa, err := db.Query("SELECT nome_completo, data_nasc, sexo, telefone, bairro, rua, numero, complemento, homem, etilista, tabagista, lesao_bucal FROM pacientes ORDER BY nome_completo")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	defer pesquisa.Close()
	var armazenamento []Pacientes
	for pesquisa.Next() {
		armazenar := Pacientes{}
		err := pesquisa.Scan(&armazenar.Nome, &armazenar.DataNasc, &armazenar.Sexo, &armazenar.Telefone, &armazenar.Bairro, &armazenar.Rua, &armazenar.Numero, &armazenar.Complemento, &armazenar.Homem, &armazenar.Etilista, &armazenar.Tabagista, &armazenar.LesaoBucal)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		quebrar := strings.Split(armazenar.DataNasc, "/")
		if armazenar.Complemento != "" {
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro + "," + armazenar.Complemento
		} else {
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro
		}
		if armazenar.LesaoBucal == "Sim" {
			if armazenar.Homem == "Sim" {
				if armazenar.Etilista == "Sim" && armazenar.Tabagista == "Sim" {
					armazenar.Fatores = "Homem/Etilista/Tabagista/Feridas Bucais"
				} else if armazenar.Etilista == "Não" && armazenar.Tabagista == "Sim" {
					armazenar.Fatores = "Homem/Tabagista/Feridas Bucais"
				} else if armazenar.Etilista == "Sim" && armazenar.Tabagista == "Não" {
					armazenar.Fatores = "Homem/Etilista/Feridas Bucais"
				}
			}
			if armazenar.Homem == "Não" {
				if armazenar.Etilista == "Sim" && armazenar.Tabagista == "Sim" {
					armazenar.Fatores = "Mulher/Etilista/Tabagista/Feridas Bucais"
				} else if armazenar.Etilista == "Não" && armazenar.Tabagista == "Sim" {
					armazenar.Fatores = "Mulher/Tabagista/Feridas Bucais"
				} else if armazenar.Etilista == "Sim" && armazenar.Tabagista == "Não" {
					armazenar.Fatores = "Mulher/Etilista/Feridas Bucais"
				}
			}
			now := time.Now()
			dia, _ := strconv.Atoi(quebrar[0])
			mes, _ := strconv.Atoi(quebrar[1])
			ano, _ := strconv.Atoi(quebrar[2])
			armazenar.Idade = now.Year() - ano
			if int(now.Month()) < mes || (int(now.Month()) == mes && now.Day() < dia) {
				armazenar.Idade--
			}
			armazenar.Usuario = usuarioLogin
			armazenar.Primeira = primeiraletraLogin
			armazenamento = append(armazenamento, armazenar)
		}
	}
	var armazenadoPgAlto []Pacientes
	armazenar2 := Pacientes{}
	armazenar2.Usuario = usuarioLogin
	armazenar2.Primeira = primeiraletraLogin
	armazenadoPgAlto = append(armazenadoPgAlto, armazenar2)
	for _, armazenado := range armazenamento {
		if masculino == "Masculino" && armazenado.Sexo == "Masculino" {
			armazenado.TemDados = true
			if feridasbucais == "FeridasBucais" && tabagista == "Tabagista" && etilista == "Etilista" && armazenado.LesaoBucal == "Sim" && armazenado.Tabagista == "Sim" && armazenado.Etilista == "Sim" {
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			}
			if feridasbucais == "FeridasBucais" && tabagista == "Tabagista" && etilista == "" && armazenado.LesaoBucal == "Sim" && armazenado.Tabagista == "Sim" && armazenado.Etilista == "Não" {
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			}
			if feridasbucais == "FeridasBucais" && tabagista == "" && etilista == "Etilista" && armazenado.LesaoBucal == "Sim" && armazenado.Tabagista == "Não" && armazenado.Etilista == "Sim" {
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			}
			if feridasbucais == "FeridasBucais" && tabagista == "" && etilista == "" && armazenado.LesaoBucal == "Sim" && armazenado.Tabagista == "Não" && armazenado.Etilista == "Não" {
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			}
			if feridasbucais == "" && tabagista == "" && etilista == "" && armazenado.LesaoBucal == "Sim"{
				armazenado.TemDados = true
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			}
		}
		if feminino == "Feminino" && armazenado.Sexo == "Feminino" {
			armazenado.TemDados = true
			if feridasbucais == "FeridasBucais" && tabagista == "Tabagista" && etilista == "Etilista" && armazenado.LesaoBucal == "Sim" && armazenado.Tabagista == "Sim" && armazenado.Etilista == "Sim" {
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			}
			if feridasbucais == "FeridasBucais" && tabagista == "Tabagista" && etilista == "" && armazenado.LesaoBucal == "Sim" && armazenado.Tabagista == "Sim" && armazenado.Etilista == "Não" {
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			}
			if feridasbucais == "FeridasBucais" && tabagista == "" && etilista == "Etilista" && armazenado.LesaoBucal == "Sim" && armazenado.Tabagista == "Não" && armazenado.Etilista == "Sim" {
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			}
			if feridasbucais == "FeridasBucais" && tabagista == "" && etilista == "" && armazenado.LesaoBucal == "Sim" && armazenado.Tabagista == "Não" && armazenado.Etilista == "Não" {
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			}
			if feridasbucais == "" && tabagista == "" && etilista == "" && armazenado.LesaoBucal == "Sim"{
				armazenado.TemDados = true
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			}
		}
		if feminino == "" && masculino == "" {
			if feridasbucais == "FeridasBucais" && tabagista == "Tabagista" && etilista == "Etilista" && armazenado.LesaoBucal == "Sim" && armazenado.Tabagista == "Sim" && armazenado.Etilista == "Sim" {
				armazenado.TemDados = true
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			}
			if feridasbucais == "FeridasBucais" && tabagista == "Tabagista" && etilista == "" && armazenado.LesaoBucal == "Sim" && armazenado.Tabagista == "Sim" && armazenado.Etilista == "Não" {
				armazenado.TemDados = true
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			}
			if feridasbucais == "FeridasBucais" && tabagista == "" && etilista == "Etilista" && armazenado.LesaoBucal == "Sim" && armazenado.Tabagista == "Não" && armazenado.Etilista == "Sim" {
				armazenado.TemDados = true
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			}
			if feridasbucais == "FeridasBucais" && tabagista == "" && etilista == "" && armazenado.LesaoBucal == "Sim" && armazenado.Tabagista == "Não" && armazenado.Etilista == "Não" {
				armazenado.TemDados = true
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			}
			if feridasbucais == "" && tabagista == "" && etilista == "" && armazenado.LesaoBucal == "Sim"{
				armazenado.TemDados = true
				if idade1 == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade2 == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade3 == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade4 == "70+" && armazenado.Idade > 70 {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
				if idade1 == "" && idade2 == "" && idade3 == "" && idade4 == "" {
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			}
		}
	}

	err = templates.ExecuteTemplate(w, "pg-alto.html", armazenadoPgAlto)
	if err != nil {
		return
	}
}

func executarPgAltoFiltroPorNome(w http.ResponseWriter, r *http.Request) {
	busca := strings.TrimSpace(r.FormValue("nome"))
	buscar, err := db.Query("SELECT nome_completo, data_nasc, sexo, telefone, bairro, rua, numero, complemento, homem, etilista, tabagista, lesao_bucal FROM pacientes WHERE nome_completo LIKE concat('%', text($1), '%')", busca)
	if err != nil {
		fmt.Println(err)
	}
	defer buscar.Close()
	var armazenamento []Pacientes
	armazenar2 := Pacientes{}
	armazenar2.Usuario = usuarioLogin
	armazenar2.Primeira = primeiraletraLogin
	armazenamento = append(armazenamento, armazenar2)
	for buscar.Next() {
		armazenar := Pacientes{}
		err = buscar.Scan(&armazenar.Nome, &armazenar.DataNasc, &armazenar.Sexo, &armazenar.Telefone, &armazenar.Bairro, &armazenar.Rua, &armazenar.Numero, &armazenar.Complemento, &armazenar.Homem, &armazenar.Etilista, &armazenar.Tabagista, &armazenar.LesaoBucal)
		if err != nil {
			panic(err.Error())
		}
		quebrar := strings.Split(armazenar.DataNasc, "/")
		if armazenar.Complemento != "" {
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro + "," + armazenar.Complemento
		} else {
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro
		}
		if armazenar.LesaoBucal == "Sim" {
			if armazenar.Homem == "Sim" {
				if armazenar.Etilista == "Sim" && armazenar.Tabagista == "Sim" {
					armazenar.Fatores = "Homem/Etilista/Tabagista/Feridas Bucais"
				} else if armazenar.Etilista == "Não" && armazenar.Tabagista == "Sim" {
					armazenar.Fatores = "Homem/Tabagista/Feridas Bucais"
				} else if armazenar.Etilista == "Sim" && armazenar.Tabagista == "Não" {
					armazenar.Fatores = "Homem/Etilista/Feridas Bucais"
				}
			}
			if armazenar.Homem == "Não" {
				if armazenar.Etilista == "Sim" && armazenar.Tabagista == "Sim" {
					armazenar.Fatores = "Mulher/Etilista/Tabagista/Feridas Bucais"
				} else if armazenar.Etilista == "Não" && armazenar.Tabagista == "Sim" {
					armazenar.Fatores = "Mulher/Tabagista/Feridas Bucais"
				} else if armazenar.Etilista == "Sim" && armazenar.Tabagista == "Não" {
					armazenar.Fatores = "Mulher/Etilista/Feridas Bucais"
				}
			}
			now := time.Now()
			dia, _ := strconv.Atoi(quebrar[0])
			mes, _ := strconv.Atoi(quebrar[1])
			ano, _ := strconv.Atoi(quebrar[2])
			armazenar.Idade = now.Year() - ano
			if int(now.Month()) < mes || (int(now.Month()) == mes && now.Day() < dia) {
				armazenar.Idade--
			}
			armazenar.TemDados = true
			armazenar.Usuario = usuarioLogin
			armazenar.Primeira = primeiraletraLogin
			armazenamento = append(armazenamento, armazenar)
		}
	}
	err = templates.ExecuteTemplate(w, "pg-alto.html", armazenamento)
	if err != nil {
		return
	}
}

func executarPgAbsenteista(w http.ResponseWriter, _ *http.Request) {
	u := UsuarioNoDashboard{Usuario: usuarioLogin, Primeira: primeiraletraLogin}
	err := templates.ExecuteTemplate(w, "pg-absenteista.html", u)
	if err != nil {
		return
	}
}

func executarFormPreenchido(w http.ResponseWriter, r *http.Request) {
	passarNome := &nomePaciente
	nome := r.FormValue("Nome")
	*passarNome = nome
	risco := r.FormValue("Risco")
	pesquisa, err := db.Query("SELECT * FROM pacientes")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	defer pesquisa.Close()
	var armazenamento []PacienteFormularioPreenchido
	for pesquisa.Next() {
		armazenar := PacienteFormularioPreenchido{}
		err = pesquisa.Scan(&armazenar.ID, &armazenar.Nome, &armazenar.DataNasc, &armazenar.CPF, &armazenar.NomeMae, &armazenar.Sexo, &armazenar.CartaoSus, &armazenar.Telefone, &armazenar.Email, &armazenar.CEP, &armazenar.Bairro, &armazenar.Rua, &armazenar.Numero, &armazenar.Complemento, &armazenar.Homem, &armazenar.Etilista, &armazenar.Tabagista, &armazenar.LesaoBucal, &armazenar.DataCadastro, &armazenar.UltimaVisita)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		armazenamento = append(armazenamento, armazenar)
	}
	for _, armazenado := range armazenamento {
		if armazenado.Nome == nome {
			if risco == "Baixo" {
				armazenado.BaixoRisco = true
			} else if risco == "Medio" {
				armazenado.MedioRisco = true
			} else {
				armazenado.AltoRisco = true
			}
			if armazenado.Etilista == "Sim" {
				armazenado.IsEtilista = true
			}
			if armazenado.Tabagista == "Sim" {
				armazenado.IsTabagista = true
			}
			if armazenado.Homem == "Sim" {
				armazenado.IsHomem = true
			}
			if armazenado.LesaoBucal == "Sim" {
				armazenado.IsLesaoBucal = true
			}
			now := time.Now()
			ultimavisita := strings.Split(armazenado.UltimaVisita, "-")
			armazenado.UltimaVisita = ultimavisita[2] + "/" + ultimavisita[1] + "/" + ultimavisita[0]
			datacadastro := strings.Split(armazenado.DataCadastro, "-")
			armazenado.DataCadastro = datacadastro[2] + "/" + datacadastro[1] + "/" + datacadastro[0]
			diaVisita, _ := strconv.Atoi(ultimavisita[2])
			mesVisita, _ := strconv.Atoi(ultimavisita[1])
			if diaVisita <= now.Day() && mesVisita < int(now.Month()){
				armazenado.MaisDeUmMes = true
			} else if (int(now.Month()) - mesVisita) >= 2{
				armazenado.MaisDeUmMes = true
			}
			datanascimento := strings.Split(armazenado.DataNasc, "/")
			armazenado.DataNasc = datanascimento[0] + "/" + datanascimento[1] + "/" + datanascimento[2]
			cnsq := strings.Split(Cns, "")
			cboq := strings.Split(Cbo, "")
			cnesq := strings.Split(Cnes, "")
			ineq := strings.Split(Ine, "")
			cbo1 := cboq[0]
			cbo2 := cboq[1]
			cbo3 := cboq[2]
			cbo4 := cboq[3]
			cbo5 := cboq[4]
			cbo6 := cboq[5]
			primeironome := strings.Split(nome, " ")
			armazenado.PrimeiroNome = primeironome[0]
			armazenado.CNS = cnsq
			armazenado.CNES = cnesq
			armazenado.CBO1 = cbo1
			armazenado.CBO2 = cbo2
			armazenado.CBO3 = cbo3
			armazenado.CBO4 = cbo4
			armazenado.CBO5 = cbo5
			armazenado.CBO6 = cbo6
			armazenado.INE = ineq
			armazenado.Usuario = usuarioLogin
			armazenado.PrimeiraLetra = primeiraletraLogin
			err = templates.ExecuteTemplate(w, "formpreenchido.html", armazenado)
			if err != nil {
				return
			}
		}
	}
}

func executarFormPreenchidoVindoDoMaps(w http.ResponseWriter, r *http.Request) {
	passarNome := &nomePaciente
	cep := r.FormValue("cep")
	risco := r.FormValue("risco")
	nome := r.FormValue("nome")
	*passarNome = nome
	pesquisa, err := db.Query("SELECT * FROM pacientes WHERE cep=$1", cep)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	defer pesquisa.Close()
	var armazenamento []PacienteFormularioPreenchido
	for pesquisa.Next() {
		armazenar := PacienteFormularioPreenchido{}
		err = pesquisa.Scan(&armazenar.ID, &armazenar.Nome, &armazenar.DataNasc, &armazenar.CPF, &armazenar.NomeMae, &armazenar.Sexo, &armazenar.CartaoSus, &armazenar.Telefone, &armazenar.Email, &armazenar.CEP, &armazenar.Bairro, &armazenar.Rua, &armazenar.Numero, &armazenar.Complemento, &armazenar.Homem, &armazenar.Etilista, &armazenar.Tabagista, &armazenar.LesaoBucal, &armazenar.DataCadastro, &armazenar.UltimaVisita)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		armazenamento = append(armazenamento, armazenar)
	}
	for _, armazenado := range armazenamento {
		if risco == "Baixo" {
			armazenado.BaixoRisco = true
		} else if risco == "Medio" {
			armazenado.MedioRisco = true
		} else {
			armazenado.AltoRisco = true
		}
		if armazenado.Etilista == "Sim" {
			armazenado.IsEtilista = true
		}
		if armazenado.Tabagista == "Sim" {
			armazenado.IsTabagista = true
		}
		if armazenado.Homem == "Sim" {
			armazenado.IsHomem = true
		}
		if armazenado.LesaoBucal == "Sim" {
			armazenado.IsLesaoBucal = true
		}
		now := time.Now()
		ultimavisita := strings.Split(armazenado.UltimaVisita, "-")
		armazenado.UltimaVisita = ultimavisita[2] + "/" + ultimavisita[1] + "/" + ultimavisita[0]
		datacadastro := strings.Split(armazenado.DataCadastro, "-")
		armazenado.DataCadastro = datacadastro[2] + "/" + datacadastro[1] + "/" + datacadastro[0]
		diaVisita, _ := strconv.Atoi(ultimavisita[2])
		mesVisita, _ := strconv.Atoi(ultimavisita[1])
		if diaVisita <= now.Day() && mesVisita < int(now.Month()){
			armazenado.MaisDeUmMes = true
		} else if (int(now.Month()) - mesVisita) >= 2{
			armazenado.MaisDeUmMes = true
		}
		datanascimento := strings.Split(armazenado.DataNasc, "/")
		armazenado.DataNasc = datanascimento[0] + "/" + datanascimento[1] + "/" + datanascimento[2]
		cnsq := strings.Split(Cns, "")
		cboq := strings.Split(Cbo, "")
		cnesq := strings.Split(Cnes, "")
		ineq := strings.Split(Ine, "")
		cbo1 := cboq[0]
		cbo2 := cboq[1]
		cbo3 := cboq[2]
		cbo4 := cboq[3]
		cbo5 := cboq[4]
		cbo6 := cboq[5]
		primeironome := strings.Split(nome, " ")
		armazenado.PrimeiroNome = primeironome[0]
		armazenado.CNS = cnsq
		armazenado.CNES = cnesq
		armazenado.CBO1 = cbo1
		armazenado.CBO2 = cbo2
		armazenado.CBO3 = cbo3
		armazenado.CBO4 = cbo4
		armazenado.CBO5 = cbo5
		armazenado.CBO6 = cbo6
		armazenado.INE = ineq
		armazenado.Usuario = usuarioLogin
		armazenado.PrimeiraLetra = primeiraletraLogin
		err = templates.ExecuteTemplate(w, "formpreenchido.html", armazenado)
		if err != nil {
			return
		}
	}
}

func alterarDataUltimaVisitaFormPreenchido(w http.ResponseWriter, r *http.Request) {
	novaDataCadastro := r.FormValue("novaVisita")
	nome := r.FormValue("nome")
	risco := r.FormValue("risco")
	cpf := r.FormValue("cpf")
	_, err := db.Exec(`UPDATE pacientes SET ultima_visita=$1 WHERE cpf=$2`, novaDataCadastro, cpf)
	if err != nil {
		return
	}
	pesquisa, err := db.Query("SELECT * FROM pacientes WHERE cpf=$1", cpf)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer pesquisa.Close()
	var armazenamento []PacienteFormularioPreenchido
	for pesquisa.Next() {
		armazenar := PacienteFormularioPreenchido{}
		err = pesquisa.Scan(&armazenar.ID, &armazenar.Nome, &armazenar.DataNasc, &armazenar.CPF, &armazenar.NomeMae, &armazenar.Sexo, &armazenar.CartaoSus, &armazenar.Telefone, &armazenar.Email, &armazenar.CEP, &armazenar.Bairro, &armazenar.Rua, &armazenar.Numero, &armazenar.Complemento, &armazenar.Homem, &armazenar.Etilista, &armazenar.Tabagista, &armazenar.LesaoBucal, &armazenar.DataCadastro, &armazenar.UltimaVisita)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		armazenamento = append(armazenamento, armazenar)
	}
	for _, armazenado := range armazenamento {
		if armazenado.Nome == nome {
			if risco == "baixo" {
				armazenado.BaixoRisco = true
			} else if risco == "medio" {
				armazenado.MedioRisco = true
			} else {
				armazenado.AltoRisco = true
			}
			if armazenado.Etilista == "Sim" {
				armazenado.IsEtilista = true
			}
			if armazenado.Tabagista == "Sim" {
				armazenado.IsTabagista = true
			}
			if armazenado.Homem == "Sim" {
				armazenado.IsHomem = true
			}
			if armazenado.LesaoBucal == "Sim" {
				armazenado.IsLesaoBucal = true
			}
			now := time.Now()
			armazenarQtdVisitadosMaisDeUmMes := &qtdVisitadosMaisDeUmMes
			ultimavisita := strings.Split(armazenado.UltimaVisita, "-")
			armazenado.UltimaVisita = ultimavisita[2] + "/" + ultimavisita[1] + "/" + ultimavisita[0]
			diaVisita, _ := strconv.Atoi(ultimavisita[2])
			mesVisita, _ := strconv.Atoi(ultimavisita[1])
			if diaVisita <= now.Day() && mesVisita < int(now.Month()){
				armazenado.MaisDeUmMes = true
			} else if (int(now.Month()) - mesVisita) >= 2{
				armazenado.MaisDeUmMes = true
			} else{
				*armazenarQtdVisitadosMaisDeUmMes--
			}
			datacadastro := strings.Split(armazenado.DataCadastro, "-")
			armazenado.DataCadastro = datacadastro[2] + "/" + datacadastro[1] + "/" + datacadastro[0]
			datanascimento := strings.Split(armazenado.DataNasc, "/")
			armazenado.DataNasc = datanascimento[0] + "/" + datanascimento[1] + "/" + datanascimento[2]
			cnsq := strings.Split(Cns, "")
			cboq := strings.Split(Cbo, "")
			cnesq := strings.Split(Cnes, "")
			ineq := strings.Split(Ine, "")
			cbo1 := cboq[0]
			cbo2 := cboq[1]
			cbo3 := cboq[2]
			cbo4 := cboq[3]
			cbo5 := cboq[4]
			cbo6 := cboq[5]
			primeironome := strings.Split(nome, " ")
			armazenado.PrimeiroNome = primeironome[0]
			armazenado.CNS = cnsq
			armazenado.CNES = cnesq
			armazenado.CBO1 = cbo1
			armazenado.CBO2 = cbo2
			armazenado.CBO3 = cbo3
			armazenado.CBO4 = cbo4
			armazenado.CBO5 = cbo5
			armazenado.CBO6 = cbo6
			armazenado.INE = ineq
			armazenado.Usuario = usuarioLogin
			armazenado.PrimeiraLetra = primeiraletraLogin
			err = templates.ExecuteTemplate(w, "formpreenchido.html", armazenado)
			if err != nil {
				return
			}
		}
	}
}

func generatePDF(w http.ResponseWriter, r *http.Request) {
	pesquisa, err := db.Query("SELECT * FROM pacientes WHERE nome_completo=$1", nomePaciente)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer pesquisa.Close()

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 22)
	pdf.SetMargins(10.0, 10.0, 10.0)

	for pesquisa.Next() {

		var armazenar PacienteFormularioPreenchido
		err := pesquisa.Scan(&armazenar.ID, &armazenar.Nome, &armazenar.DataNasc, &armazenar.CPF, &armazenar.NomeMae, &armazenar.Sexo, &armazenar.CartaoSus, &armazenar.Telefone, &armazenar.Email, &armazenar.CEP, &armazenar.Bairro, &armazenar.Rua, &armazenar.Numero, &armazenar.Complemento, &armazenar.Homem, &armazenar.Etilista, &armazenar.Tabagista, &armazenar.LesaoBucal, &armazenar.DataCadastro, &armazenar.UltimaVisita)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
		quebrar := strings.Split(armazenar.DataCadastro, "-")
		armazenar.DataCadastro = quebrar[2] + "/" + quebrar[1] + "/" + quebrar[0]
		imageFile:="imagenspg/logo projeto 2.png"
		infoPtr:= pdf.RegisterImage(imageFile, "")
		pageWidth, _ := pdf.GetPageSize()
		y:= (-infoPtr.Width()+ pageWidth)/2
		pdf.ImageOptions("imagenspg/logo projeto 2.png", y, 10, 40, 0, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "") 
		pdf.Ln(40)
		texto:="SOBREVIDAS ACS"
		width:= pdf.GetStringWidth(texto)
		x:=width-pageWidth
		pdf.SetX(x)
		pdf.SetTextColor(0,140,255)
		pdf.SetDrawColor(0, 50, 140)
		pdf.CellFormat(width-15.0, 9.0, texto, "B", 1, "C", false, 0, "")
		pdf.Ln(30)
		x=10.0
		pdf.SetTextColor(0,0,0)
		pdf.SetX(x)
		pdf.SetDrawColor(30, 30, 80)
		azulcell(pdf,"Dados do Paciente")
		comumcell(pdf, "Nome Completo", 50.0, 0, "B")
		comumcell(pdf, armazenar.Nome, 140.0, 1, "")
		comumcell(pdf, "Sexo", 30.0, 0, "B")
		comumcell(pdf, armazenar.Sexo, 50.0, 0, "")
		comumcell(pdf, "Data de Nascimento", 60.0, 0, "B")
		comumcell(pdf, armazenar.DataNasc, 50.0, 1, "")
		comumcell(pdf, "Nome da Mãe", 50.0, 0, "B")
		comumcell(pdf, armazenar.NomeMae, 140.0, 1, "")
		comumcell(pdf, "CPF", 30.0, 0, "B")
		comumcell(pdf, armazenar.CPF, 50.0, 0, "")
		comumcell(pdf, "Cartão SUS", 40, 0, "B")
		comumcell(pdf, armazenar.CartaoSus, 70, 1, "")
		pdf.Ln(6)
		azulcell(pdf,"Dados de Contato")
		comumcell(pdf, "Telefone", 30.0, 0, "B")
		comumcell(pdf, armazenar.Telefone, 50.0, 0, "")
		comumcell(pdf,"Email", 20.0, 0, "B")
		comumcell(pdf, armazenar.Email, 90.0, 1, "")
		pdf.Ln(6)
		azulcell(pdf,"Dados Geográficos")
		comumcell(pdf, "CEP",30.0, 0, "B")
		comumcell(pdf, armazenar.CEP, 50.0, 0, "")
		comumcell(pdf, "Bairro", 25.0, 0, "B")
		comumcell(pdf, armazenar.Bairro, 85.0, 1, "")
		comumcell(pdf, "Rua", 20.0, 0, "B")
		comumcell(pdf, armazenar.Rua, 40.0, 0, "")
		comumcell(pdf, "Número", 30.0, 0, "B")
		comumcell(pdf, armazenar.Numero, 15.0, 0, "")
		comumcell(pdf, "Complemento", 40.0, 0, "B")
		comumcell(pdf, armazenar.Complemento, 45.0, 1, "")
		pdf.Ln(6)
		azulcell(pdf,"Fatores de Prioridade")
		prioridade(pdf,"Homem (idade igual ou superior a 40 anos)",armazenar.Homem)
		prioridade(pdf,"Etilista (consome álcool regularmente)",armazenar.Etilista)
		prioridade(pdf,"Tabagista (consome cigarro regularmente)",armazenar.Tabagista)
		prioridade(pdf,"Com lesão bucal suspeita",armazenar.LesaoBucal)
		pdf.Ln(4)
		comumcell(pdf, "Data de Cadastro", 60.0,0, "B")
		comumcell(pdf, armazenar.DataCadastro, 130.0,1, "")
	err = pesquisa.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=output.pdf")

	err = pdf.Output(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
}

func azulcell(pdf *gofpdf.Fpdf, texto string,){
	tr := pdf.UnicodeTranslatorFromDescriptor("")
		pdf.SetFont("Arial", "B", 16)
		pdf.SetFillColor(0, 125, 230)
		pdf.SetTextColor(255, 255, 255)
		pdf.CellFormat(190, 9.0, tr(texto), "", 1, "L", true, 0, "")
		pdf.Ln(2)

		//titulos
}

func comumcell(pdf *gofpdf.Fpdf, texto string, width float64, pula int, negrito string){
	tr := pdf.UnicodeTranslatorFromDescriptor("")
		pdf.SetFont("Arial", negrito, 15)
		pdf.SetFillColor(235, 240, 255)
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(width, 9.0, tr(texto), "1", pula, "L", true, 0, "")

		// celula padrao
}

func prioridade(pdf *gofpdf.Fpdf, texto string, condicao string) {
	tr := pdf.UnicodeTranslatorFromDescriptor("")
		pdf.SetDrawColor(10, 10, 10)
		pdf.SetFillColor(255,255,255)
		pdf.SetFont("Arial", "B", 14)
		pdf.SetTextColor(0,0,0)
		pdf.CellFormat(130, 9.0, tr(texto), "1", 0, "C", true, 0, "")
		pdf.SetFont("Arial", "", 14)

	pdf.SetFillColor(170, 255, 170)
	mensagem := "Negativo"
	if condicao == "Sim" {
		pdf.SetFillColor(255, 170, 170)
		mensagem = "Positivo"
	}
	pdf.SetFont("Arial", "", 16)
		pdf.CellFormat(60, 9.0, tr(mensagem), "1", 1, "C", true, 0, "")
	
	
	//mudando as corzinha
}
