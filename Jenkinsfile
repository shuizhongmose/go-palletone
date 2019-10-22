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
                    script {
                        env.BUILD_STATUS = 'failed'
                    }

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