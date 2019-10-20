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
        APPLICATION_DIR='application'
        USER_CONTRACT_DIR = 'usercontract'
        BLACKLIST_DIR='blacklist'
        GO111MODULE = 'on'
        LOG_NAME='log.html'
        REPORT_NAME='report.html'

        IS_RUN_UT='false'
        IS_RUN_USER_CONTRACT='false'
        IS_RUN_DIGITAL = 'false'
        IS_RUN_APPLICATION = 'true'
        IS_RUN_MEDIATOR_VOTE = 'false'

        IS_RUN_DEPOSIT = 'true'
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

        IS_UPLOAD = 'false'
    }
    stages {
        stage('UT') {
            when {
                environment name: 'IS_RUN_UT', value: 'true'
            }
            steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    sh 'export PATH=${GOPATH}:${PATH}'
                    sh 'go build -mod=vendor ./cmd/gptn'
                    sh 'make gptn'
                    sh 'go test -mod=vendor ./...'
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

                        chmod +x upload.sh
                        ./upload.sh
                    '''

                    sh 'pkill gptn'
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

                        chmod +x upload.sh
                        ./upload.sh
                    '''

                    sh 'pkill gptn'
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
                        	sh 'ls ${BASE_DIR}/bdd/node'
                            sh 'netstat -ap | grep gptn'

                        	script {
                        	    if (env.IS_RUN_DEPOSIT == 'true') {
                        	        sh '''
                        	            cd ${BASE_DIR}/bdd/dct
                        	            pwd
                        	            ls ../node
                                        ./deposit_test.sh 7
                                    '''
                        	    }
                        	}
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
                        		cd ${BASE_DIR}/bdd
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
                        	sh '''
                        		cd ${BASE_DIR}
                        		zip -j ./bdd/logs/oneNode_log.zip ./bdd/node/log/*
                        		zip -j ./bdd/logs/gasToken_log.zip ./bdd/GasToken/node/log/*
                        		./bdd/upload2Ftp.sh ${FTP_PWD} ${TRAVIS_BRANCH} ${TRAVIS_BUILD_NUMBER}
                        	'''
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
                        	sh '''
                        		cd ${BASE_DIR}
                        		zip -j ./bdd/logs/zMulti-node.zip ./logs/zMulti-node/*
                        		./bdd/upload2Ftp.sh ${FTP_PWD} ${TRAVIS_BRANCH} ${TRAVIS_BUILD_NUMBER}
                        		cd ${BASE_DIR}/bdd
                        		source ./targz_node.sh
                        		./upload2Ftp.sh ${FTP_PWD} ${TRAVIS_BRANCH} ${TRAVIS_BUILD_NUMBER}
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
            stages{
                stage('Running') {
                    steps {
                        catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                        	sh '''
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
                        		zip -j ./bdd/logs/application_log.zip ./bdd/application/node/log/*
                        		./bdd/upload2Ftp.sh ${FTP_PWD} ${TRAVIS_BRANCH} ${TRAVIS_BUILD_NUMBER}
                        	'''
                        }
                    }
                }
            }
        }
    }
}