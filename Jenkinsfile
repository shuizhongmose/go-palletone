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
            }
            post {
                failure {
                    echo '1111111'
                    script { env.BUILD_STATUS = 'failed' }
                }
                success {
                    echo '22222'
                }
            }
        }
    }
    post {
        always {
            echo env.BUILD_STATUS
        }
    }
}