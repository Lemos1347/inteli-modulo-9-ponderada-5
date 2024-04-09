# Integração do Simulador com Kafka e HiveMQ

## 1. Objetivo

Este projeto visa criar um simulador de dispositivos IoT que utiliza o protocolo MQTT para enviar informações simuladas baseadas em dados de sensores reais, nesse caso o [Sensor de Radiação Solar Sem Fio HOBOnet RXW-LIB-900](https://sigmasensors.com.br/produtos/sensor-de-radiacao-solar-sem-fio-hobonet-rxw-lib-900), integrar o HiveMQ com o Confluent Kafka (provedor cloud de kafka) e consumir um tópico kafka.

## 2. Funcionamento

https://github.com/Lemos1347/inteli-modulo-9-ponderada-5/assets/99190347/75959021-1469-4f7d-b168-1a49294fdaee

## 3. Como Instalar e Rodar

### Pré-requisitos

> [!IMPORTANT]
> - Credenciais de um broker (HiveMQ)[https://www.hivemq.com/]
> - Credenciais de um (Confluent Kafka)[https://www.confluent.io/cloud-kafka/?utm_medium=sem&utm_source=google&utm_campaign=ch.sem_br.brand_tp.prs_tgt.confluent-brand_mt.xct_rgn.latam_lng.eng_dv.all_con.confluent-kafka-general&utm_term=confluent%20kafka&creative=&device=c&placement=&gad_source=1&gclid=Cj0KCQjwq86wBhDiARIsAJhuphkrRGOZ_lDlLOz8C3bWHd4Cykwwn26xRC2PAxkOTi1efuqHSQ8vVRYaAi7LEALw_wcB]
> - (Docker)[https://www.docker.com/] e (docker-compose)[https://docs.docker.com/compose/] instalado
> - [Golang](https://go.dev/doc/install) em sua máquina.

### Instalação

Clone o repositório para a sua máquina local:

```bash
git clone https://github.com/Lemos1347/inteli-modulo-9-ponderada-5
cd inteli-modulo-9-ponderada-5
```

Assegure as instalações das libs:
```bash
go mod tidy
```

### Execução

> [!IMPORTANT]
> Antes de rodar a aplicacao, garanta que na pasta `configs` tenha o arquivo `client.properties` com as credenciais do seu Confluent Kafka e os seguinte arquivos dotenv com as suas respectivas credenciais:
> `.env.subscriber`
> ```bash
> DATABASE_PORT=3002
> DATABASE_HOST=localhost
> DATABASE_USER=user
> DATABASE_PASSWORD=password
> DATABASE_NAME=Ponderada-4
> PROPERTIES_FILE_PATH=../../configs/client.properties
> ```
> `.env.publisher`
> ```bash
> DATABASE_PORT=3002
> DATABASE_HOST=localhost
> DATABASE_USER=user
> DATABASE_PASSWORD=password
> DATABASE_NAME=Ponderada-4
> BROKER_URL=YOUR_BROKER_URL
> BROKER_USER=ponderada5-pub
> BROKER_PASSWORD=Ponderada5-nicola-pub
> CSV_PATH=./assets/data/dados_sensor_radiacao_solar.csv
>```

#### Publisher

O publisher nada mais e que a simulacao do sensor solar, ele ira publicar dados simulados em um topico MQTT: `ponderada5/sensors`. Para rodar ele, e apenas necessario executar o binario, passar o path para o arquivo do `.env` e a quantidade de sensores que voce quer criar (caso você não passe nada, ele emulará os sensores que já existem, caso nao exista nenhum sensor e voce não passe nenhum valor ele não rodará):

```bash
go run cmd/publisher/publisher.go ./configs/.env.publisher <numero>
```

Pronto! A partir de agora os dados serão publicados no topico "ponderada5/sensors".

#### Subscriber

O subscriber nada mais é que um servidor que organizará os dados advindos do topico kafka `ponderada5` e os armazenará no banco de dados.

Para rodar ele, você precisa executar o binario e passar o path para o arquivo do `.env`:

```bash
go run cmd/subscriber/subscriber.go ./configs/.env.subscriber
```

Pronto! Agora os dados enviados serão armazenados no banco de dados!

## 3. Estrutura do Projeto

O projeto é composto por:

- `build`: os executáveis e o docker da aplicacao
- `cmd`: entrypoint de todos os códigos em golang
- `configs`: arquivos de configuração do projeto
- `internal`: codigos nao exportáveis do pacote golang

## 4. Testes 
https://github.com/Lemos1347/inteli-modulo-9-ponderada-5/assets/99190347/d9ede3b5-cfda-4702-8bb5-c71125e47a0f

Todos os testes podem ser encontrados na pasta [`tests`](./tests/). Para rodar os testes, basta executar o comando na root do projeto:

```bash
go test -v -cover ./...
```

> [!IMPORTANT]
> Para rodar os testes, é necessário ter o arquivo `client.properties` com as credenciais do seu Confluent Kafka o seguitne arquivo dotenv na pasta `configs`:
> `.env.test`
> ```bash
> BROKER_URL=YOUR_BROKER_URL
> BROKER_USER=ponderada5-test-pub
> BROKER_PASSWORD=ponderada5-test-pubPassword
> PROPERTIES_FILE_PATH=../configs/client.properties
> ```