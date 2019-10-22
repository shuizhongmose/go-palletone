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
                    echo '11111'
                    script {
                        BUILD_STATUS='failed'
                    }
                    echo 'BUILD_STATUS is now ${BUILD_STATUS}'
                }
                success {
                    echo '22222'
                }
            }
        }
    }
    post {
        always {
            echo 'BUILD_STATUS is now ${BUILD_STATUS}'
            script {
                if (env.BUILD_STATUS=='failed') {
                    sh 'exit 1'
                }
            }
        }
    }
}