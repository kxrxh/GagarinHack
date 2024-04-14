<script setup>
import { AuthService } from '@/services/AuthService';
import UILabeledInput from '@/components/ui/UILabeledInput.vue';
import UIButton from '@/components/ui/UIButton.vue';
</script>

<template>
    <section class="bg-gray-50 dark:bg-gray-900 h-screen">
        <div class="flex flex-col items-center justify-center px-6 py-8 mx-auto md:h-screen lg:py-0">
            <div class="w-full bg-white rounded-lg shadow dark:border md:mt-0 sm:max-w-md xl:p-0 dark:bg-gray-800 dark:border-gray-700">
                <div class="p-6 space-y-4 md:space-y-6 sm:p-8">
                    <h1 class="text-xl font-bold leading-tight tracking-tight text-gray-900 md:text-2xl dark:text-white text-center">
                        Авторизация
                    </h1>
                    <div class="space-y-4 md:space-y-6" action="#">
                        <UILabeledInput
                            v-model="login"
                            type="text"
                            property="login"
                            placeholder="">
                            Почта
                        </UILabeledInput>
                        <UILabeledInput
                            v-model="password"
                            type="password"
                            property="password"
                            placeholder="">
                            Пароль
                        </UILabeledInput>
                        <UIButton @click="auth" :disabled="loading || !isInputsSet" classExtension="w-full px-5 py-2.5">
                            Войти
                        </UIButton>
                    </div>
                </div>
            </div>
        </div>
    </section>
</template>

<script>
import { useWebApp } from 'vue-tg';
export default {
    name: 'AuthView',
    data() {
        return {
            login: "",
            password: "",

            loading: false
        }
    },
    computed: {
        isInputsSet() {
            return this.login.length > 0 && this.password.length > 0;
        }
    },
    methods: {
        auth() {
            if(this.login.length < 1) {
                this.$notify({text:"Введите почту или номер телефона", type: "error"});
                return;
            }
            if(this.password.length < 1) {
                this.$notify({text:"Введите пароль", type: "error"});
                return;
            }
            AuthService.login(this.login, this.password, (data) => {
                this.$notify({text:"Успешная авторизация", type: "success"});
                let token = data.access_token;
                this.$cookies.set('token', token, "30d");
                sessionStorage.setItem("is_auth", true);
                useWebApp().sendData(JSON.stringify(
                        {
                            type: "auth",
                            data: {
                                token: this.$cookies.get("token"),
                                email: this.login
                            }
                        }
                    ));
                if(sessionStorage.getItem("fallback")) {
                    this.$router.push({path: sessionStorage.getItem("fallback")});
                } else useWebApp().close();
                return;
            }, (error) => {
                Object.values(error.response.data.errors).flat().forEach(message => {
                    this.$notify({text: message, type: "error"});
                });
            })
        }
    }
}
</script>