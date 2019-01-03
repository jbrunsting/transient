<template>
  <div class="login">
    <form @submit.prevent="login">
      <input placeholder="username" v-model="username">
      <input type="password" placeholder="password" v-model="password">
      <button type="submit">{{ submitText }}</button>
    </form>
    <Error class="error empty">
      Please enter your username and password
    </Error>
    <Error class="error login">
      Username or password incorrect. <a href="todo">Forgot password?</a>
    </Error>
    <Error class="error unknown">
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
    props: {
        submitText: {
            type: String,
            default: 'Login',
        },
        apiPath: {
            type: String,
            default: '/api/user/login',
        },
    },
    components: {
        Error,
    },
    methods: {
        login() {
            /* eslint-disable no-param-reassign */
            this.$el.querySelectorAll('.error').forEach((c) => {
                c.style.visibility = 'hidden';
            });

            if (this.username === '' || this.password === '') {
                this.$el.querySelector('.empty.error').style.visibility = 'visible';
                return;
            }

            const identification = {
                username: this.username,
                password: this.password,
            };

            this.$http.post(this.apiPath, identification)
                .then(() => {
                    this.$emit('login');
                }).catch((e) => {
                    if (e.response.status === 401) {
                        this.$el.querySelector('.login.error').style.visibility = 'visible';
                    } else {
                        console.log(e.response);
                        this.$el.querySelector('.unknown.error').style.visibility = 'visible';
                        if (e) {
                            console.log(`${JSON.stringify(e)}`);
                        }
                    }
                });
            /* eslint-enable no-param-reassign */
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

.error {
  position: absolute;
  visibility: hidden;
}
</style>
