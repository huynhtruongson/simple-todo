Feature: Login

    Background:
        Given user has already created an account

    Scenario Outline: Login successfully
        When user login with "valid credential" by "<loginType>"
        Then user login successfully by "<loginType>"
        Examples:
            | loginType |
            # | grpc      |
            | http      |

    Scenario: Login failed
        When user login with "<errType>" by "<loginType>"
        Then user will get error message "<errMsg>"
        Examples:
            | errType            | errMsg                                        | loginType |
            | empty username     | Login error: Username or Password is required | grpc      |
            | empty password     | Login error: Username or Password is required | grpc      |
            | incorrect username | Login error: Username or Password incorrect   | grpc      |
            | incorrect password | Login error: Username or Password incorrect   | grpc      |
            | empty username     | Username or Password is required              | htpp      |
            | empty password     | Username or Password is required              | htpp      |
            | incorrect username | Username or Password incorrect                | htpp      |
            | incorrect password | Username or Password incorrect                | htpp      |

