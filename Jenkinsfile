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
                script {
                    if (currentBuild.result != 0) {
                        env.BUILD_STATUS = 'failure'
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