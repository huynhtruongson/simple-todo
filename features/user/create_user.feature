Feature: Create user

    Scenario: Create user successfully
        When user creates account with "valid information"
        Then the account is created successfully
        And user can login with username and password

    Scenario Outline: Create user failed
        When user creates account with "<errType>"
        Then user will get error message "<errMsg>"
        Examples:
            | errType           | errMsg                                    |
            | empty fullname    | Fullname is empty                         |
            | empty username    | Username is empty                         |
            | invalid username  | Username length is less than 6 characters |
            | existing username | Username has already existed              |
            | empty email       | Email is empty                            |
            | invalid email     | Email is invalid                          |
            | existing email    | Email has already existed                 |
            | empty password    | Password is empty                         |
            | invalid password  | Password length is less than 6 characters |