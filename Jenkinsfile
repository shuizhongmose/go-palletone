pipeline {
    agent any
    options {
        disableConcurrentBuilds()
        checkoutToSubdirectory('/home/JGitlabGo/src/github.com/palletone/go-palletone')
        gitLabConnection('palletone-gitlab')
    }
    environment {
        GOPATH = '/home/JGitlabGo'
        GOCACHE = '/home/JGitlabGo/caches/go'

        BASE_DIR = '/home/JGitlabGo/src/github.com/palletone/go-palletone'
        ALL_LOG_PATH = '/home/JGitlabGo/src/github.com/palletone/go-palletone/bdd/node/log/all.log'
        GAS_TOKEN_ALL_LOG_PATH = '/home/JGitlabGo/src/github.com/palletone/go-palletone/bdd/GasToken/node/log/all.log'
        BDD_LOG_PATH = '/home/JGitlabGo/src/github.com/palletone/go-palletone/bdd/logs'
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
        APPLICATION_DIR='application'
        USER_CONTRACT_DIR = 'usercontract'
        BLACKLIST_DIR='blacklist'
        GO111MODULE = 'on'
        LOG_NAME='log.html'
        REPORT_NAME='report.html'

        IS_RUN_UT='true'
        IS_RUN_USER_CONTRACT='false'
        IS_RUN_DIGITAL = 'false'
        IS_RUN_APPLICATION = 'false'
        IS_RUN_MEDIATOR_VOTE = 'false'

        IS_RUN_DEPOSIT = 'false'
        IS_RUN_TESTCONTRACTCASES = 'false'
        IS_RUN_CREATE_TRANS = 'false'
        IS_RUN_20CONTRACT = 'false'
        IS_RUN_721SEQENCE = 'false'
        IS_RUN_721UDID = 'false'
        IS_RUN_GASTOKEN = 'false'
        IS_RUN_VOTE = 'false'
        IS_RUN_MULTIPLE = 'false'
        IS_RUN_LIGHT = 'false'
        IS_RUN_BLACKLIST = 'false'

        IS_UPLOAD = 'true'
    }
    stages {
        stage('Install Requirements') {
            steps{
                sh '''
                    pip install --upgrade pip
                    pip install robotframework==2.8.5
                    pip install requests
                    pip install robotframework-requests
                    pip install demjson
                    pip install pexpect
                    apt-get install expect
                    apt-get install lftp
                '''
            }
        }
        stage('UT') {
            when {
                environment name: 'IS_RUN_UT', value: 'true'
            }
            steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    sh '''
                        export PATH=${GOPATH}:${PATH}
                        cd ${BASE_DIR}
                        go build -mod=vendor ./cmd/gptn
                        make gptn
                        go test -mod=vendor ./...
                    '''
                }
            }
        }
        stage('User Contract BDD') {
            when {
                environment name: 'IS_RUN_USER_CONTRACT', value: 'true'
            }
            steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    sh '''
                        cd ${BASE_DIR}/bdd/UserContract/scripts
                        chmod +x start.sh
                        ./start.sh
                        pkill gptn
                        sleep 2
                    '''
                    script {
                        if (env.IS_UPLOAD=='true') {
                            sh '''
                                cd ${BASE_DIR}/bdd/UserContract/scripts
                                chmod +x upload.sh
                                ./upload.sh
                            '''
                        }
                    }
                }
            }
        }
        stage('Digital Identity BDD') {
            when {
                environment name: 'IS_RUN_DIGITAL', value: 'true'
            }
            steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    sh '''
                        cd ${BASE_DIR}/bdd/Digital-Identity/scripts
                        chmod +x start.sh
                        ./start.sh
                        pkill gptn
                        sleep 2
                    '''
                    script {
                        if (env.IS_UPLOAD=='true') {
                            sh '''
                                cd ${BASE_DIR}/bdd/Digital-Identity/scripts
                                chmod +x upload.sh
                                ./upload.sh
                            '''
                        }
                    }
                }
            }
        }
        stage('One Node BDD') {
			stages {
                stage('Build') {
                    steps {
                        catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                            sh '''
                                cd ${BASE_DIR}
                        		rm -rf bdd/GasToken/node
                        		rm -rf bdd/node/palletone

                                export GO111MODULE=on
                        		go build -mod=vendor ./cmd/gptn
                        		cp gptn bdd/node/
                        		mkdir bdd/GasToken/node
                        		cp gptn bdd/GasToken/node
                        		cd bdd/node
                        		chmod +x gptn
                        		python init.py
                        		nohup ./gptn &
                        		sleep 15
                        		ls ./
                        		pwd
                        	'''
                            sh 'netstat -ap | grep gptn'
                        }
                    }
                }
                stage('Deposit') {
                    when {
                        environment name: 'IS_RUN_DEPOSIT', value: 'true'
                    }
                    steps {
                        catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                            sh '''
                                cd ${BASE_DIR}/bdd/dct
                                ./deposit_test.sh 7
                            '''
                        }
                    }
                }
                stage('Blacklist') {
                    when {
                        environment name: 'IS_RUN_BLACKLIST', value: 'true'
                    }
                    steps {
                        catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                            sh '''
                                cd ${BASE_DIR}/bdd/blacklist
                                ./blacklist_test.sh
                            '''
                        }
                    }
                }
                stage('ContractTestcases') {
                    when {
                        environment name: 'IS_RUN_TESTCONTRACTCASES', value: 'true'
                    }
                    steps {
                        catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                            sh '''
                                cd ${BASE_DIR}/bdd/contract/testcases
                                chmod +x ./test_start.sh
                                ./test_start.sh
                            '''
                        }

                    }
                }
                stage('Create Transaction') {
                    when {
                        environment name: 'IS_RUN_CREATE_TRANS', value: 'true'
                    }
                    steps {
                        catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                            sh '''
                                cd ${BASE_DIR}/bdd
                                python -m robot.run -d ${BDD_LOG_PATH}/${CREATE_TRANS_DIR} -i normal ./testcase/createTrans
                            '''
                        }
                    }
                }
                stage('PRC720 Contract') {
                    when {
                        environment name: 'IS_RUN_20CONTRACT', value: 'true'
                    }
                    steps {
                        catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                            sh '''
                                cd ${BASE_DIR}/bdd
                                python -m robot.run -d ${BDD_LOG_PATH}/${CONTRACT20_DIR} -i normal ./testcase/crt20Contract
                            '''
                        }
                    }
                }
                stage('PRC721 Contract') {
                    when {
                        environment name: 'IS_RUN_721SEQENCE', value: 'true'
                    }
                    steps {
                        catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                            sh '''
                                cd ${BASE_DIR}/bdd
                                python -m robot.run -d ${BDD_LOG_PATH}/${SEQENCE721_DIR} -i normal ./testcase/crt721Seqence
                            '''
                        }
                    }
                }
                stage('PRC721 UDID') {
                    when {
                        environment name: 'IS_RUN_721UDID', value: 'true'
                    }
                    steps {
                        catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                            sh '''
                                cd ${BASE_DIR}/bdd
                                python -m robot.run -d ${BDD_LOG_PATH}/${UDID721_DIR} -i normal ./testcase/crt721UDID
                            '''
                        }
                    }
                }
                stage('Vote') {
                    when {
                        environment name: 'IS_RUN_VOTE', value: 'true'
                    }
                    steps {
                        catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                        	sh '''
                        		cd ${BASE_DIR}/bdd
                        		python -m robot.run -d ${BDD_LOG_PATH}/${VOTECONTRACT_DIR} -i normal ./testcase/voteContract
                        	'''
                        }
                    }
                }
                stage('Gas Token') {
                    when {
                        environment name: 'IS_RUN_GASTOKEN', value: 'true'
                    }
                    steps {
                        catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                        	sh '''
                        		cd ${BASE_DIR}/bdd/GasToken
                        		chmod +x ./init_gas_token.sh
                        		./init_gas_token.sh
                        		sleep 15
                        		python -m robot.run -d ${BDD_LOG_PATH}/${GAS_TOKEN_DIR} ./testcases
                        	'''
                        }
                    }
                }
                stage('After Running') {
                    steps {
                        sh '''
                            killall gptn
                            sleep 2
                        '''
                    }
                }
                stage('Upload Logs') {
                    when {
                        environment name: 'IS_UPLOAD', value: 'true'
                    }
                    steps {
                        catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                            script {
                                if (env.IS_RUN_GASTOKEN=='true') {
                                    sh 'zip -j ./bdd/logs/gasToken_log.zip ./bdd/GasToken/node/log/*'
                                }
                            }

                            withCredentials([usernamePassword(credentialsId: 'UPLOAD_TO_FTP_ID', passwordVariable: 'FTP_PWD', usernameVariable: 'FTP_USNAME')]) {
                                sh '''
                                    cd ${BASE_DIR}
                                    zip -j ./bdd/logs/oneNode_log.zip ./bdd/node/log/*
                                    echo ${FTP_PWD}
                                    chmod +x ./bdd/uploadJenkins2Ftp.sh
                                    ./bdd/uploadJenkins2Ftp.sh ${FTP_PWD} ${BRANCH_NAME} ${BUILD_NUMBER}
                                '''
                            }
                        }
                    }
                }
			}
        }
        stage('Multiple Nodes BDD') {
            stages {
                stage('Running') {
                    steps {
                        catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                            sh '''
                                cd ${BASE_DIR}
                                export GO111MODULE=on
                                make gptn
                                cp build/bin/gptn bdd/node
                                cd bdd/node
                                chmod -R +x *
                                sudo -H chmod +w /etc/hosts
                                sudo -H sed -i 's/127.0.0.1 localhost/127.0.0.1/g' /etc/hosts
                                sudo -H sed -i '$a0.0.0.0 localhost' /etc/hosts
                                ./launchMultipleNodes.sh
                                netstat -ap | grep gptn
                                grep "mediator_interval" node1/ptn-genesis.json
                                grep "maintenance_skip_slots" node1/ptn-genesis.json
                                cd ${BASE_DIR}/bdd
                                mkdir -p ${BDD_LOG_PATH}
                            '''
                        }
                    }
                }
                stage('Run Multiple') {
                    when {
                        environment name: 'IS_RUN_MULTIPLE', value: 'true'
                    }
                    steps {
                        catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                        	sh '''
                        		cd ${BASE_DIR}
                        		python -m robot.run -d ${BDD_LOG_PATH}/${MULTIPLE_DIR} -i normal ./testcase/zMulti-node
                        	'''
                        }
                    }
                }
                stage('Run Light') {
                    when {
                        environment name: 'IS_RUN_LIGHT', value: 'true'
                    }
                    steps {
                        catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                        	sh '''
                        		cd ${BASE_DIR}/light
                        		chmod +x ./bddstart.sh
                        		./bddstart.sh
                        	'''
                        }
                    }
                }
                stage('After Running') {
                    steps {
                        catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                            sh '''
                                killall gptn
                                sleep 2
                            '''
                        }

                    }
                }
                stage('Upload Logs') {
                    when {
                        environment name: 'IS_UPLOAD', value: 'true'
                    }
                    steps {
                        catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                        	sh '''
                        		cd ${BASE_DIR}
                        		zip -j ./bdd/logs/zMulti-node.zip ./logs/zMulti-node/*
                        		chmod +x ./bdd/uploadJenkins2Ftp.sh
                        		./bdd/uploadJenkins2Ftp.sh ${FTP_PWD} ${BRANCH_NAME} ${BUILD_NUMBER}
                        		cd ${BASE_DIR}/bdd
                        		source ./targz_node.sh
                        		./uploadJenkins2Ftp.sh ${FTP_PWD} ${BRANCH_NAME} ${BUILD_NUMBER}
                        	'''
                        }
                    }
                }
            }
        }
        stage('Application BDD') {
            when {
                environment name: 'IS_RUN_APPLICATION', value: 'true'
            }
            steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    sh '''
                        cd ${BASE_DIR}
                        export GO111MODULE=on
                        go build -mod=vendor ./cmd/gptn
                        mkdir bdd/application/node
                        cp gptn bdd/application/node
                        cd ./bdd/application
                        chmod +x ./init.sh
                        ./init.sh
                        sleep 15
                        python -m robot.run -d ${BDD_LOG_PATH}/${APPLICATION_DIR} .

                        killall gptn
                        sleep 2
                    '''

                    script {
                        if (env.IS_UPLOAD == 'true') {
                            sh '''
                                cd ${BASE_DIR}
                                zip -j ./bdd/logs/application_log.zip ./bdd/application/node/log/*
                        		chmod +x ./bdd/uploadJenkins2Ftp.sh
                                ./bdd/uploadJenkins2Ftp.sh ${FTP_PWD} ${BRANCH_NAME} ${BUILD_NUMBER}
                            '''
                        }
                    }
                }
            }
        }
    }
    post {
        success {
          updateGitlabCommitStatus name: 'Jenkins CI Integration', state: 'success'
        }
        failure {
          updateGitlabCommitStatus name: 'Jenkins CI Integration', state: 'failed'
        }
    }
}
