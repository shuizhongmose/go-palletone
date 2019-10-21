pipeline {
    agent any
    environment {
        BUILD_STATUS = 'success'
    }
    stages {
        stage('TEST') {
            steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    sh 'exit 1'
                }
                scripts {
                    if (currentBuild.result != 0) {
                        env.BUILD_STATUS = 'failure'
                    }
                }
            }
        }
        post {
            always {
                sh 'echo env.BUILD_STATUS'
            }
        }
    }
}