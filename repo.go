package main

import "fmt"

var currentId int

var todos Todos
var config ConfigData

// Give us some seed data
func init() {
    RepoCreateTodo(Todo{Name: "Write presentation"})
    RepoCreateTodo(Todo{Name: "Host meetup"})
    RepoCreateUpdateConfig(ConfigData{DifficultyRating: 2, Startup: StartupData{AiNationCount: 1, StartupCash: 20, AiStartupCash: 10, AiAggressiveness: 1, StartupIndependentTown: 15, StartupRawSite: 5, DifficultyLevel: 1 }})
}

func RepoFindTodo(id int) Todo {
    for _, t := range todos {
        if t.Id == id {
            return t
        }
    }
    // return empty Todo if not found
    return Todo{}
}

func RepoCreateTodo(t Todo) Todo {
    currentId += 1
    t.Id = currentId
    todos = append(todos, t)
    return t
}

func RepoDestroyTodo(id int) error {
    for i, t := range todos {
        if t.Id == id {
            todos = append(todos[:i], todos[i+1:]...)
            return nil
        }
    }
    return fmt.Errorf("Could not find Todo with id of %d to delete", id)
}

func RepoFindConfig() ConfigData {
    return config
}

func RepoCreateUpdateConfig(c ConfigData) ConfigData {
    config = c
    return config
}


