<template>
    <div class="test">
        <span v-html="titleHtml"></span>
        <h1>{{title}}</h1>
        <input type="text" v-model="title"><br/>
        <p v-if="showName">{{user.firstName}}</p>
        <p v-text="user.firstName" v-if="showName"></p>
        <p v-else="showName">Not Signed in</p>
        <ul>
            <li v-for="item in items">{{item.title}}</li>
        </ul>
        <button v-on:click="greet('Hello World')">Say Greeting</button>
        <input type="text" v-on:keyup="pressKey" v-on:keyup.enter="enterHit">
    </div>
</template>

<script>
    export default {
      name: 'test',
      data () {
        return {
          title: 'Hello World',
          titleHtml: '<h1>Hello World</h1>',
          user: {
            firstName: 'John',
            lastName: 'Harrison'
          },
          showName: false,
          items: [
              {title: 'Item 1'},
              {title: 'Item 2'},
              {title: 'Item 3'},
              {title: 'Item 4'}
          ]
        }
      },
      methods: {
        greet: function (greeting) {
          alert(greeting)
        },
        pressKey: function (event) {
          console.log(event.target.value)
        },
        enterHit: function () {
          console.log('You hit enter')
        }
      }
    }
</script>

<style scoped>

</style>
