pipeline {
    agent none
    stages {
        stage('One') {
            steps{
                echo 'this is one'
            }
        }
        stage('Sequential') {
            stages {
                stage('In Sequential 1') {
                    steps {
                        echo "In Sequential 1"
                    }
                }
                stage('In Sequential 2') {
                    steps {
                        echo "In Sequential 2"
                    }
                }
                stage('Parallel In Sequential') {
                    parallel {
                        stage('In Parallel 1') {
                            steps {
                                echo "In Parallel 1"
                            }
                        }
                        stage('In Parallel 2') {
                            steps {
                                echo "In Parallel 2"
                            }
                        }
                    }
                }
            }
        }
        stage('Two') {
            steps {
                echo 'this is two'
            }
        }
    }
}