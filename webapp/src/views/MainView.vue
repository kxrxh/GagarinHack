<script setup>
import UILabeledInput from '../components/ui/UILabeledInput.vue'
import UIButton from '../components/ui/UIButton.vue'
import '../assets/css/loader.css'
</script>
<template>
    <section class="bg-gray-50 dark:bg-gray-900 h-screen">
        <div class="flex flex-col items-center justify-center px-6 py-8 mx-auto md:h-screen lg:py-0 animate animate-fade animate-ease-in-out animate-duration-250 animate-once">
            <div class="w-full bg-white rounded-lg shadow dark:border md:mt-0 sm:max-w-md xl:p-0 dark:bg-gray-800 dark:border-gray-700">
                <div class="p-6 space-y-4 md:space-y-6 sm:p-8">
                    <h1 class="text-xl font-bold leading-tight tracking-tight text-gray-900 md:text-2xl dark:text-white text-center">Код памяти</h1>
                    <div v-if="block == 0" class="space-y-4 md:space-y-6 animate animate-fade animate-ease-in-out animate-duration-250 animate-once">
                      <UILabeledInput
                        type="text"
                        class="animate animate-fade animate-ease-in-out animate-duration-350 animate-once"
                        v-model="answers[index]">
                        {{ questions[index] }}
                      </UILabeledInput>
                      <UIButton classExtension="w-full py-2.5" @click="next" :disabled="answers[index].length < 1">
                        Продолжить
                      </UIButton>
                      <p class="text-sm font-light text-gray-500 dark:text-gray-400 text-center !mt-4">
                          <a href="#" class="font-medium text-primary-600 hover:underline dark:text-primary-500" @click="back" :disabled="index == 0">
                            Вернуться на предыдущий
                          </a><br><br>
                          <a href="#" class="font-medium text-primary-600 hover:underline dark:text-primary-500" @click="skip">
                            Пропустить вопрос
                          </a>
                      </p>
                    </div>
                    <div v-if="block == 1" class="space-y-4 md:space-y-6 animate animate-fade animate-ease-in-out animate-duration-250 animate-once text-white text-center">
                      Создание кода памяти...<br>
                      <span class="loader"></span>
                    </div>
                    <div v-if="block == 2" class="space-y-2 md:space-y-4 animate animate-fade animate-ease-in-out animate-duration-250 animate-once text-white text-center">
                      <div class="grid gap-2 mb-2 grid-cols-2">
                        <UIButton color="warning" classExtension="py-2.5" @click="regenerate">
                          Перегенерировать
                        </UIButton>
                        <UIButton color="danger" classExtension="py-2.5" @click="reset">
                          Сбросить
                        </UIButton>
                      </div>
                      <UILabeledInput
                        :textarea="true"
                        type="text"
                        class="animate animate-fade animate-ease-in-out animate-duration-350 animate-once"
                        v-model="results[0]">
                        Предложенный код памяти №1:
                      </UILabeledInput>
                      <UIButton classExtension="w-full py-2.5" @click="finish">
                        Принять первый код памяти
                      </UIButton>
                      <UILabeledInput
                        :textarea="true"
                        type="text"
                        class="animate animate-fade animate-ease-in-out animate-duration-350 animate-once"
                        v-model="results[1]">
                        Предложенный код памяти №2:
                      </UILabeledInput>
                      <UIButton classExtension="w-full py-2.5" @click="finish">
                        Принять второй код памяти
                      </UIButton>
                    </div>
                </div>
            </div>
        </div>
    </section>
</template>

<script>
export default {
    name: "MainView",
    data() {
      return {
        index: 0,
        questions: ["Как его зовут?", "В каком городе он родился?", "Что вас в нем радовало?"],
        answers: ["", "", ""],
        results: ["", ""],
        block: 2
      }
    },
    methods: {
      next() {
        this.index++;
        this.checkEnd();
      },
      skip() {
        this.index++;
        // TODO recreate question
        this.checkEnd();
      },
      back() {
        if(this.index == 0) return;
        this.index--;
      },
      checkEnd() {
        if(this.index >= this.questions.length) {
          this.generate();
        }
      },
      finish() {
        // TODO send to bot
      },
      regenerate() {
        this.generate();
      },
      generate() {
        this.block = 1;
        // TODO send to backend
      },
      reset() {
        this.block = 0;
        // TODO remove questions and create new
      }
    },
    mounted() {

    }
}
</script>