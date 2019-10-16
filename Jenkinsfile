pipeline {
    agent any
    options {
        disableConcurrentBuilds()
        checkoutToSubdirectory('/home/JGithubgo/src/github.com/palletone/go-palletone')
    }
    environment {
        GOPATH = '/home/JGithubgo'
        GOCACHE = '/home/JGithubgo/caches/go'

        BASE_DIR = '/home/JGithubgo/src/github.com/palletone/go-palletone'
        ALL_LOG_PATH = '/home/JGithubgo/src/github.com/palletone/go-palletone/bdd/node/log/all.log'
        GAS_TOKEN_ALL_LOG_PATH = '/home/JGithubgo/src/github.com/palletone/go-palletone/bdd/GasToken/node/log/all.log'
        BDD_LOG_PATH = '/home/JGithubgo/src/github.com/palletone/go-palletone/bdd/logs'
        CREATE_TRANS_DIR = 'createTrans'
        CONTRACT20_DIR = 'crt20Contract'
        SEQENCE721_DIR = 'crt721Seqence'
        UDID721_DIR = 'crt721UDID'
        VOTECONTRACT_DIR = 'voteContract'
        MULTIPLE_DIR = 'zMulti-node'
        DIGITAL_IDENTITY_DIR = 'Digital-Identitycert'
        DEPOSIT_DIR = 'deposit'
        GAS_TOKEN_DIR = 'gasToken'
        MEDIATOR_VOTE_DIR = 'meidatorvote'
        USER_CONTRACT_DIR = 'usercontract'
        GO111MODULE = 'on'
        FTP_PWD = 'Pallet2018'
    }
    stages {
        stage('Install Requirements') {
            steps{
                echo 'hello world'
            }
        }
        try {
            stage('UT') {
                steps {
                    catchError {
                        sh 'export PATH=${GOPATH}:${PATH}'
                        sh 'cd ${BASE_DIR}'
                        sh 'go build -mod=vendor ./cmd/gptn'
                        sh 'make gptn'
                        sh 'go test -mod=vendor ./...'
                    }
                    echo stageResult.result
                }
            }
        } catch (Exception e) {
            echo 'Stage failed, but we continue'
        }
        stage('User Contract BDD') {
            steps {
                sh '''
                    cd ${BASE_DIR}/bdd/UserContract/scripts
                    ls
                    chmod +x start.sh
                    ./start.sh

                    chmod +x upload.sh
                    ./upload.sh
                '''

                sh 'pkill gptn'
            }
        }
        stage('Digital Identity BDD') {
            steps {
                sh '''
                    cd ${BASE_DIR}/bdd/Digital-Identity/scripts
                    chmod +x start.sh
                    ./start.sh

                    chmod +x upload.sh
                    ./upload.sh
                '''

                sh 'pkill gptn'
            }
        }
    }
}