<template>
  <div class="login">
    <form @submit.prevent="login">
      <input placeholder="username" v-model="username">
      <input type="password" placeholder="password" v-model="password">
      <button type="submit">Login</button>
    </form>
    <Error class="login-error">
      Username or password incorrect. <a href="todo">Forgot password?</a>
    </Error>
    <Error class="unknown-error">
      Could not login, please try again later
    </Error>
  </div>
</template>

<script>
import Error from '@/components/Error.vue';

export default {
    name: 'Login',
    data() {
        return {
            username: '',
            password: '',
        };
    },
    components: {
        Error,
    },
    methods: {
        login() {
            this.$el.querySelectorAll('.login-error').forEach((c) => {
                c.style.visibility = 'hidden';
            });

            const identification = {
                username: this.username,
                password: this.password,
            };

            this.$http.post('/api/user/login', identification)
                .then(() => {
                    this.$emit('login');
                }).catch((e) => {
                    if (e.response.status === 401) {
                        this.$el.querySelectorAll('.login-error').forEach((c) => {
                            c.style.visibility = 'visible';
                        });
                    } else {
                        this.$el.querySelectorAll('.unknown-error').forEach((c) => {
                            c.style.visibility = 'visible';
                        });
                        console.log(`${JSON.stringify(e.response.status)}`);
                    }
                });
        },
    },
};
</script>

<style scoped lang="scss">
@import "../styles/settings.scss";

.login {
  position: relative;
}

input {
  margin-top: 0;
  margin-bottom: 0;
  width: 120px;
}

button {
  margin: 0 $margin0;
}

form {
  display: flex;
  align-items: center;
}

.login-error {
  position: absolute;
  visibility: hidden;
}

.unknown-error {
  position: absolute;
  visibility: hidden
}
</style>
