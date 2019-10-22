pipeline {
    agent any
    environment {
        BUILD_STATUS = 'success'
    }
    stages {
        stage('TEST') {
            steps {
                catchError(buildResult: 'FAILURE', stageResult: 'FAILURE') {
                    sh 'exit 1'
                }

                script {
                    if (currentBuild.result=='FAILURE') {
                        echo '111111'
                        BUILD_STATUS = "failed"
                        echo "BUILD_STATUS is now '${BUILD_STATUS}'"
                    }
                }
            }
        }
    }
    post {
        always {
            echo "BUILD_STATUS is now '${BUILD_STATUS}'"
            script {
                if (env.BUILD_STATUS=='failed') {
                    sh 'exit 1'
                }
            }
        }
    }
}