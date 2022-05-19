<template>
  <div>
    <h2 v-if="errorMessage.length > 0">{{errorMessage}}</h2>
    <section class="todoapp" v-cloak>
      <header class="header">
        <h1>todos</h1>
        <div class="buttons">
          <ul class="filters">
            <li>
              <a @click="saveAs">Save As</a>
            </li>
            <li>
              <a @click="loadNewList">Load</a>
            </li>
          </ul>
        </div>
        <input
          class="new-todo"
          autofocus
          autocomplete="off"
          placeholder="What needs to be done?"
          v-model="newTodo"
          @keyup.enter="addTodo"
        >
      </header>
      <section class="main" v-show="todos.length">
        <ul class="todo-list">
          <li
            class="todo"
            v-for="todo in todos"
            :key="todo.id"
            :class="{completed: todo.completed, editing: todo == editedTodo}"
          >
            <div class="view">
              <input class="toggle" type="checkbox" v-model="todo.completed">
              <label @dblclick="editTodo(todo)">{{todo.title}}</label>
              <button class="destroy" @click="removeTodo(todo)"></button>
            </div>
            <input
              class="edit"
              type="text"
              v-model.lazy="todo.title"
              @keyup.enter="doneEdit(todo)"
              @blur="doneEdit(todo)"
              @keyup.esc="cancelEdit(todo)"
            >
          </li>
        </ul>
      </section>
    </section>
  </div>
</template>

<script setup>
import "./assets/css/base.css";
import "./assets/css/app.css";

import { EventsOnMultiple, LogInfo } from "../wailsjs/runtime"
import * as backend from "../wailsjs/go/main/Todos"
import { onMounted, reactive, ref, watch } from "vue";

const todos = reactive([])
const loading = ref(false)
const errorMessage = ref('')
const editedTodo = ref(null)
const newTodo = ref('')

const addTodo = function() {
  var value = newTodo.value && newTodo.value.trim();
  if (!value) {
    return;
  }
  todos.push({
    id: todos.length,
    title: value,
    completed: false
  });
  newTodo.value = "";
}
 const removeTodo = function(todo) {
  var index = todos.indexOf(todo);
  todos.splice(index, 1);

  for (var i = 0; i < todos.length; i++) {
    todos[i].id = i;
  }
}

const editTodo = function(todo) {
  beforeEditCache = todo.title;
  editedTodo.value = todo;
}

const doneEdit = function(todo) {
  if (!editedTodo.value) {
    return;
  }
  editedTodo.value = null;
  todo.title = todo.title.trim();
  if (!todo.title) {
    removeTodo(todo);
  }
}
const cancelEdit = function(todo) {
  editedTodo.value = null;
  todo.title = beforeEditCache;
}
const saveAs = function() {
  backend.SaveAs(JSON.stringify(todos, null, 2));
}
const loadNewList = function() {
  backend.LoadNewList();
}
const setErrorMessage = function(message) {
  errorMessage.value = message;
  setTimeout(() => {
    errorMessage.value = "";
  }, 3000);
}
const loadList = function() {
  backend.LoadList()
    .then(list => {
      try {
        console.log('loadList result:', list);
        let todoList = JSON.parse(list);
        loading.value = true;
        todos.splice(0, todos.length)
        todos.push(...todoList)
      } catch (e) {
        setErrorMessage("Unable to load todo list");
      }
    })
    .catch(error => {
      setErrorMessage(error.message);
    });
}

watch(todos, function(todos) {
  if (loading.value) {
    loading.value = false;
    return;
  }
  backend.SaveList(JSON.stringify(todos, null, 2));
})

onMounted(()=>{
  console.log('mounted')
  LogInfo("in onMounted")
  EventsOnMultiple("filemodified", () => {
    console.log('filemodified')
    alert('filemodified')
    loadList();
  });

  EventsOnMultiple("error", (message, number) => {
    let result = number * 2;
    setErrorMessage(`${message}: ${result}`);
  });

  // Load the list at the start
  loadList();
})

</script>

<style>
h2 {
  text-align: center;
  color: white;
  background-color: red;
  min-width: 230px;
  max-width: 550px;
  padding: 1rem;
  border-radius: 0.5rem;
}

.buttons {
  height: 20px;
  padding: 10px 20px;
  box-shadow: inset 0 -2px 1px rgba(0, 0, 0, 0.1);
  text-align: center;
  border-color: rgba(175, 47, 47, 0.2);
}

.buttons ul li a {
  margin: 10px;
}

.buttons li {
  border-color: rgba(175, 47, 47, 0.1);
}

.filters li a {
  color: inherit;
  margin: 3px;
  padding: 3px 7px;
  text-decoration: none;
  border: 1px solid transparent;
  border-radius: 3px;
  border-color: rgba(100, 100, 100, 0.1);
}
.filters li a:hover {
  border-color: rgba(255, 47, 47, 0.3);
  cursor: pointer;
}
</style>