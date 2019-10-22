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
                echo currentBuild.result
                echo currentStage.result
                script {
                    if (currentBuild.result != 0) {
                        env.BUILD_STATUS = 'failed'
                    }
                    if (currentBuild.result != '0') {
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