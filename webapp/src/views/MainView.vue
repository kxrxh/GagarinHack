<script setup>
import UILabel from '@/components/ui/UILabel.vue';
import UILabeledInput from '../components/ui/UILabeledInput.vue';
import UIDropdown from '../components/ui/UIDropdown.vue';
import UIButton from '../components/ui/UIButton.vue';
import '../assets/css/loader.css';
import '@vuepic/vue-datepicker/dist/main.css';
import VueDatePicker from '@vuepic/vue-datepicker';
</script>
<template>
    <section class="bg-gray-50 dark:bg-gray-900 h-screen">
        <div class="flex flex-col items-center justify-center px-6 py-8 mx-auto md:h-screen lg:py-0 animate animate-fade animate-ease-in-out animate-duration-250 animate-once">
            <div class="w-full bg-white rounded-lg shadow dark:border md:mt-0 sm:max-w-md xl:p-0 dark:bg-gray-800 dark:border-gray-700">
                <div class="p-6 space-y-4 md:space-y-6 sm:p-8">
                    <h1 class="text-xl font-bold leading-tight tracking-tight text-gray-900 md:text-2xl dark:text-white text-center">Страница памяти</h1>
                    <div v-if="stage == STAGE.SETUP" class="space-y-4 md:space-y-6 animate animate-fade animate-ease-in-out animate-duration-250 animate-once">
                      <div>
                        <UILabel>
                          Выберите дату рождения
                        </UILabel>
                        <VueDatePicker v-model="dates[0]"
                          :enable-time-picker="false"
                          :max-date="dates[1]"
                          :format="formatter"
                          locale="ru-RU"
                          select-text="Выбрать"
                          cancel-text="Закрыть"
                          :dark="true"
                          :clearable="false"
                          class="!mt-0"
                          input-class-name="dp-custom-input"></VueDatePicker>
                      </div>
                      <div>
                        <UILabel>
                          Выберите дату смерти
                        </UILabel>
                        <VueDatePicker v-model="dates[1]"
                          :enable-time-picker="false"
                          :min-date="dates[0]"
                          :max-date="new Date()"
                          :format="formatter"
                          locale="ru-RU"
                          select-text="Выбрать"
                          cancel-text="Закрыть"
                          :dark="true"
                          :clearable="false"
                          class="!mt-0"
                          input-class-name="dp-custom-input"></VueDatePicker>
                      </div>
                      <div>
                        <UILabel>
                          Выберите пол
                        </UILabel>
                        <UIDropdown
                          v-model="sex"
                          :options='{"мужского": "Мужской", "женского": "Женский"}'>
                          Пол не выбран
                        </UIDropdown>
                      </div>
                      <UILabeledInput
                        type="text"
                        v-model="name">
                        Введите ФИО
                      </UILabeledInput>
                      <UIButton classExtension="w-full py-2.5" @click="next" :disabled="name.length < 1">
                        Продолжить
                      </UIButton>
                    </div>
                    <div v-if="stage == STAGE.AWAIT_QUESTION" class="space-y-4 md:space-y-6 animate animate-fade animate-ease-in-out animate-duration-250 animate-once text-white text-center">
                      Создание вопросов...<br>
                      <span class="loader"></span>
                    </div>
                    <div v-if="stage == STAGE.QUESTION" class="space-y-4 md:space-y-6 animate animate-fade animate-ease-in-out animate-duration-250 animate-once">
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
                          <a href="#" class="font-medium text-primary-600 hover:underline dark:text-primary-500" @click="skip">
                            Пропустить вопрос<br><br>
                          </a>
                          <a href="#" class="font-medium text-primary-600 hover:underline dark:text-primary-500" @click="back" v-if="index != 0">
                            Вернуться на предыдущий
                          </a>
                      </p>
                    </div>
                    <div v-if="stage == STAGE.CREATION" class="space-y-4 md:space-y-6 animate animate-fade animate-ease-in-out animate-duration-250 animate-once text-white text-center">
                      Создание эпитафии...<br>
                      <span class="loader"></span>
                    </div>
                    <div v-if="stage == STAGE.VIEW" class="space-y-2 md:space-y-4 animate animate-fade animate-ease-in-out animate-duration-250 animate-once text-white text-center">
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
                        Предложенная эпитафия №1:
                      </UILabeledInput>
                      <UIButton classExtension="w-full py-2.5" @click="finish">
                        Принять первую эпитафию
                      </UIButton>
                      <UILabeledInput
                        :textarea="true"
                        type="text"
                        class="animate animate-fade animate-ease-in-out animate-duration-350 animate-once"
                        v-model="results[1]">
                        Предложенная эпитафия №2:
                      </UILabeledInput>
                      <UIButton classExtension="w-full py-2.5" @click="finish">
                        Принять вторую эпитафию
                      </UIButton>
                    </div>
                </div>
            </div>
        </div>
    </section>
</template>

<script>
import { MemoryService } from '@/services/MemoryService';
const MODEL = "gigachat";
const STAGE = {
  SETUP: 100,
  QUESTION: 101,
  CREATION: 102,
  VIEW: 103
}
export default {
    name: "MainView",
    components: {
      VueDatePicker
    },
    data() {
      return {
        index: 0,
        dates: [new Date(), new Date()],
        name: "",
        sex: "мужского",
        questions: [],
        answers: [],
        results: ["", ""],
        stage: STAGE.SETUP,
        unskip: [0]
      }
    },
    methods: {
      next() {
        if(this.stage == STAGE.SETUP) {
          this.stage = STAGE.AWAIT_QUESTION;
          MemoryService.getQuestions(MODEL, this.name, this.sex, this.dates[0], (data) => {
            this.questions = data.response;
            this.answers = Array.from({ length: this.questions.length }).fill("");
            this.stage = STAGE.QUESTION;
          }, (err) => {
            // TODO handle error
            debugger;
          })
          return;
        }
        this.index++;
        this.checkEnd();
      },
      skip() {
        this.index++;
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
        this.stage = STAGE.CREATION;
        const questions = {};
        for (let i = 0; i < this.questions.length; i++) {
            const question = this.questions[i];
            const answer = this.answers[i];
            if (answer !== "") {
              questions[question] = answer;
            }
        }
        const checkContinue = () => {
          if(this.results[0] && this.results[1]) {
            this.stage = STAGE.VIEW;
          }
        }
        MemoryService.getStory(MODEL, "epitaph", this.name, this.sex, this.dates[0], this.dates[1], questions, (data) => {
          this.results[0] = data.response;
          checkContinue();
        }, (err) => {
          // TODO handle error
          debugger;
        })
        MemoryService.getStory("yandex", "epitaph", this.name, this.sex, this.dates[0], this.dates[1], questions, (data) => {
          this.results[1] = data.response;
          checkContinue();
        }, (err) => {
          // TODO handle error
          debugger;
        })
      },
      reset() {
        this.stage = STAGE.QUESTION;
        // TODO remove questions and create new
      },
      formatter(date) {
        const day = String(date.getDate()).padStart(2, '0');
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const year = date.getFullYear();

        return `${day}.${month}.${year}`;
      }
    },
    mounted() {

    }
}
</script>

<style>
.dp-custom-input {
  box-shadow: 0 0 6px rgb(14, 165, 233);
  color: rgb(14, 165, 233);
  margin-top: 0;
}
</style>