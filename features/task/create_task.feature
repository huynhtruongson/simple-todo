Feature: Create task

    Background: 
        Given user have been signed in

    Scenario: Create task successfully
        When user creates task with "valid information"
        Then the task is created successfully

    Scenario Outline: Create user failed
        When user creates task with "<errType>"
        Then user will get error message "<errMsg>"
        Examples:
            | errType        | errMsg            |
            | empty title    | Title is empty    |
            | invalid status | Status is invalid |
