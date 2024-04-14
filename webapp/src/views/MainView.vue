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
                    <div v-if="stage == STAGE.SETUP_STARTER" class="space-y-4 md:space-y-6 animate animate-fade animate-ease-in-out animate-duration-250 animate-once">
                      <UILabeledInput
                        type="text"
                        v-model="author">
                        Введите ваше ФИО (автор эпитафии)
                      </UILabeledInput>
                      <UIButton classExtension="w-full py-2.5" @click="changeStage(STAGE.SETUP_PERSON)" :disabled="author.length < 1">
                        Продолжить
                      </UIButton>
                    </div>
                    <div v-if="stage == STAGE.SETUP_PERSON" class="space-y-4 md:space-y-6 animate animate-fade animate-ease-in-out animate-duration-250 animate-once">
                      <div>
                        <UILabel>
                          Укажите пол
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
                      <UIButton classExtension="w-full py-2.5" @click="changeStage(STAGE.SETUP_PLACES)" :disabled="name.length < 1">
                        Продолжить
                      </UIButton>
                    </div>
                    <div v-if="stage == STAGE.SETUP_PLACES" class="space-y-4 md:space-y-6 animate animate-fade animate-ease-in-out animate-duration-250 animate-once">
                      <UILabeledInput
                        type="text"
                        v-model="places[0]">
                        Введите место рождения
                      </UILabeledInput>
                      <UILabeledInput
                        type="text"
                        v-model="places[1]">
                        Введите место смерти
                      </UILabeledInput>
                      <UIButton classExtension="w-full py-2.5" @click="changeStage(STAGE.SETUP_RELATIVES)" :disabled="places[0].length < 1 || places[1].length < 1">
                        Продолжить
                      </UIButton>
                    </div>
                    <div v-if="stage == STAGE.SETUP_RELATIVES" class="space-y-4 md:space-y-6 animate animate-fade animate-ease-in-out animate-duration-250 animate-once">
                      <UILabeledInput
                        type="text"
                        v-model="partner">
                        Укажите {{ sex == "мужского" ? "супругу" : "супруга" }} (при наличии)
                      </UILabeledInput>
                      <UILabeledInput
                        type="text"
                        :textarea="true"
                        v-model="children">
                        Укажите детей (при наличии)
                      </UILabeledInput>
                      <UIButton classExtension="w-full py-2.5" @click="changeStage(STAGE.SETUP_EDUCATION)">
                        Продолжить
                      </UIButton>
                    </div>
                    <div v-if="stage == STAGE.SETUP_EDUCATION" class="space-y-4 md:space-y-6 animate animate-fade animate-ease-in-out animate-duration-250 animate-once">
                      <UILabeledInput
                        type="text"
                        v-model="citizenship">
                        Укажите гражданство
                      </UILabeledInput>
                      <UILabeledInput
                        type="text"
                        v-model="education">
                        Укажите образование
                      </UILabeledInput>
                      <UILabeledInput
                        type="text"
                        v-model="career">
                        Укажите род деятельности
                      </UILabeledInput>
                      <UILabeledInput
                        type="text"
                        :textarea="true"
                        v-model="achievments">
                        Награды, премии, достижения
                      </UILabeledInput>
                      <UIButton classExtension="w-full py-2.5" @click="next" :disabled="citizenship.length < 1">
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
                      <UIButton classExtension="w-full py-2.5" @click="acceptEpitaphy(0)">
                        Принять первую эпитафию
                      </UIButton>
                      <UILabeledInput
                        :textarea="true"
                        type="text"
                        class="animate animate-fade animate-ease-in-out animate-duration-350 animate-once"
                        v-model="results[1]">
                        Предложенная эпитафия №2:
                      </UILabeledInput>
                      <UIButton classExtension="w-full py-2.5" @click="acceptEpitaphy(1)">
                        Принять вторую эпитафию
                      </UIButton>
                    </div>
                    <div v-if="stage == STAGE.PROMPT_BIOGRAPHY" class="space-y-2 md:space-y-4 animate animate-fade animate-ease-in-out animate-duration-250 animate-once text-white text-center">
                      <UILabel>
                        Хотите ли вы создать биографию {{ name }}?
                      </UILabel>
                      <div class="grid gap-2 mb-2 grid-cols-2">
                        <UIButton color="danger" classExtension="py-2.5" @click="finish">
                          Нет
                        </UIButton>
                        <UIButton classExtension="py-2.5" @click="changeStage(STAGE.SELECT_BIOGRAPHY)">
                          Да
                        </UIButton>
                      </div>
                    </div>
                    <div v-if="stage == STAGE.SELECT_BIOGRAPHY" class="space-y-2 md:space-y-4 animate animate-fade animate-ease-in-out animate-duration-250 animate-once text-white text-center">
                      <UILabel>
                        До какого периода вы хотите создать биографию?
                      </UILabel>
                      <UIButton classExtension="py-2.5" @click="selectBiography(1)">
                        До детства/юношества (Одна часть)
                      </UIButton>
                      <UIButton classExtension="py-2.5" @click="selectBiography(2)">
                        До средних лет (Две части)
                      </UIButton>
                      <UIButton classExtension="py-2.5" @click="selectBiography(3)">
                        До старости (Три части)
                      </UIButton>
                    </div>
                    <div v-if="stage == STAGE.BIOGRAPHY_QUESTION" class="space-y-4 md:space-y-6 animate animate-fade animate-ease-in-out animate-duration-250 animate-once">
                      <UILabel>
                        {{ biography[biographyIndex].stage }}
                      </UILabel>
                      <UILabeledInput
                        type="text"
                        class="animate animate-fade animate-ease-in-out animate-duration-350 animate-once"
                        v-model="biography[biographyIndex].answers[index]">
                        {{ biography[biographyIndex].questions[index] }}
                      </UILabeledInput>
                      <UIButton classExtension="w-full py-2.5" @click="next" :disabled="biography[biographyIndex].answers[index].length < 1">
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
                    <div v-if="stage == STAGE.BIOGRAPHY_CREATION" class="space-y-4 md:space-y-6 animate animate-fade animate-ease-in-out animate-duration-250 animate-once text-white text-center">
                      Создание биографии...<br>
                      <span class="loader"></span>
                    </div>
                    <div v-if="stage == STAGE.VIEW_BIOGRAPHY" class="space-y-4 md:space-y-6 animate animate-fade animate-ease-in-out animate-duration-250 animate-once">
                        <div class="grid gap-2 mb-2 grid-cols-2">
                          <UIButton color="warning" classExtension="py-2.5" @click="generateBiography">
                            Перегенерировать
                          </UIButton>
                          <UIButton color="danger" classExtension="py-2.5" @click="changeStage(STAGE.PROMPT_BIOGRAPHY)">
                            Сбросить
                          </UIButton>
                        </div>
                        <UILabel>
                          Предложенная биография
                        </UILabel>
                        <div style="max-height: 50vh; overflow: auto;">
                          <div v-for="(item, i) in biography" :key="item.stage" :class="[item.header.length < 1 ? 'hidden' : '']">
                            <UILabeledInput
                              type="text"
                              v-model="item.header">
                              Заголовок:
                            </UILabeledInput>
                            <UILabeledInput
                              :textarea="true"
                              type="text"
                              v-model="item.text">
                              Текст:
                            </UILabeledInput>
                          </div>
                          <UILabeledInput
                            :textarea="true"
                            type="text"
                            v-model="biographyEnding">
                            Заключение:
                          </UILabeledInput>
                        </div>
                        <UIButton classExtension="w-full py-2.5" @click="acceptBiography">
                          Принять биографию
                        </UIButton>
                    </div>
                    <div v-if="stage == STAGE.APPLYING" class="space-y-4 md:space-y-6 animate animate-fade animate-ease-in-out animate-duration-250 animate-once text-white text-center">
                      Отправка страницы памяти...<br>
                      <span class="loader"></span>
                    </div>
                </div>
            </div>
        </div>
    </section>
</template>

<script>
import { MemoryService } from '@/services/MemoryService';
import { ExternalService } from '@/services/ExternalService';
const MODEL = "gigachat";
const STAGE = {
  SETUP_STARTER: 1,
  SETUP_PERSON: 2,
  SETUP_PLACES: 3,
  SETUP_RELATIVES: 4,
  SETUP_EDUCATION: 5,
  QUESTION: 101,
  CREATION: 102,
  VIEW: 103,
  PROMPT_BIOGRAPHY: 200,
  SELECT_BIOGRAPHY: 201,
  BIOGRAPHY_QUESTION: 202,
  BIOGRAPHY_CREATION: 203,
  VIEW_BIOGRAPHY: 204,
  APPLYING: 300
}
export default {
    name: "MainView",
    components: {
      VueDatePicker
    },
    data() {
      return {
        author: "",
        sex: "мужского",
        name: "",
        dates: [new Date(), new Date()],
        places: ["", ""],
        partner: "",
        children: "",
        citizenship: "Россия",
        education: "",
        career: "",
        achievments: "",

        index: 0,
        questions: [],
        answers: [],
        results: ["", ""],
        stage: STAGE.SETUP_STARTER,
        unskip: [0],

        epitaphy: "",
        biographyIndex: 0,
        maxBiography: 0,
        biography: [
          {
            stage: "Детство и юношество",
            questions: [],
            answers: [],
            header: "",
            text: ""
          },
          {
            stage: "Средние годы",
            questions: [],
            answers: [],
            header: "",
            text: ""
          },
          {
            stage: "Последние годы",
            questions: [],
            answers: [],
            header: "",
            text: ""
          }
        ],
        biographyEnding: "",
        useBiography: false
      }
    },
    methods: {
      getQuestions() {
        MemoryService.getQuestions(MODEL, this.name, this.sex, this.dates[0], (data) => {
            this.questions = data.response;
            this.answers = Array.from({ length: this.questions.length }).fill("");
            this.stage = STAGE.QUESTION;
          }, (err) => {
            this.$notify({text:"Не удалось получить вопросы, попробуйте позже.", type: "error"});
            console.log(err);
          });
      },
      next() {
        if(this.stage != STAGE.QUESTION && this.stage != STAGE.BIOGRAPHY_QUESTION) {
          this.stage = STAGE.AWAIT_QUESTION;
          this.getQuestions();
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
        if(this.maxBiography != 0) {
          if(this.index >= this.biography[this.biographyIndex].answers.length) {
            if(this.biographyIndex < this.maxBiography - 1) {
              this.biographyIndex++;
              this.generateBiographyQuestions();
            } else {
              this.generateBiography();
            }
          }
        } else if(this.index >= this.questions.length) {
          this.generate();
        }
      },
      acceptEpitaphy(index) {
        this.epitaphy = this.results[index];
        this.stage = STAGE.PROMPT_BIOGRAPHY;
      },
      regenerate() {
        this.generate();
      },
      changeStage(stage) {
        this.stage = stage;
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
        this.results = ["", ""];
        MemoryService.getEpitaph(MODEL, this.name, this.sex, this.dates[0], this.dates[1], questions, (data) => {
          this.results[0] = data.response;
          checkContinue();
        }, (err) => {
          this.$notify({text:"Не удалось сгенерировать эпитафию с использованием Gigachat, попробуйте позже.", type: "error"});
          this.results[0] = "Ошибка";
          checkContinue();
          console.log(err);
        })
        MemoryService.getEpitaph("yandex", this.name, this.sex, this.dates[0], this.dates[1], questions, (data) => {
          this.results[1] = data.response;
          checkContinue();
        }, (err) => {
          this.$notify({text:"Не удалось сгенерировать эпитафию с использованием YandexGPT, попробуйте позже.", type: "error"});
          this.results[1] = "Ошибка";
          checkContinue();
          console.log(err);
        })
      },
      reset() {
        this.stage = STAGE.AWAIT_QUESTION;
        this.getQuestions();
      },
      formatter(date) {
        const day = String(date.getDate()).padStart(2, '0');
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const year = date.getFullYear();

        return `${day}.${month}.${year}`;
      },
      selectBiography(max) {
        this.maxBiography = max;
        this.generateBiographyQuestions();
      },
      generateBiographyQuestions() {
        this.stage = STAGE.AWAIT_QUESTION;
        let type = this.biographyIndex == 0 ? "youth" :
                    (this.biographyIndex == 1 ? "middle_age" : "old_age");
        MemoryService.getBiographyQuestions(MODEL, type, this.name, this.sex, this.dates[0], (data) => {
          this.biography[this.biographyIndex].questions = data.response;
          this.biography[this.biographyIndex].answers = Array.from({ length: this.biography[this.biographyIndex].questions.length }).fill("");
          this.index = 0;
          this.stage = STAGE.BIOGRAPHY_QUESTION;
        }, err => {
          this.$notify({text:"Не удалось получить вопросы, попробуйте позже.", type: "error"});
          console.log(err);
        })
      },
      generateNextBiography(i) {
        if(i == 3) {
          MemoryService.getBiographyShort(MODEL, this.biography, this.name, this.sex, this.dates[0], this.dates[1], (data) => {
            this.biographyEnding = data.response;
            this.stage = STAGE.VIEW_BIOGRAPHY;
          }, err => {
            this.$notify({text:"Не удалось сгенерировать заключение биографии, попробуйте позже.", type: "error"});
            console.log(err);
          })
          return;
        }
        let questions = {};
        let amount = 0;
        for (let j = 0; j < this.biography[i].questions.length; j++) {
            const question = this.biography[i].questions[j];
            const answer = this.biography[i].answers[j];
            if (answer !== "") {
              questions[question] = answer;
              amount++;
            }
        }
        if(amount > 0) {
          let type = i == 0 ? "youth" :
                  (i == 1 ? "middle_age" : "old_age");
          MemoryService.getBiography(MODEL, type, this.name, this.sex, this.dates[0], this.dates[1], questions, i == 0 ? "Пустая биография" : this.biography[i - 1].text, (data) => {
            this.biography[i].header = data.header;
            this.biography[i].text = data.response;
            this.generateNextBiography(i + 1);
          }, err => {
            this.$notify({text:"Не удалось сгенерировать биографию, попробуйте позже.", type: "error"});
            console.log(err);
          })
        } else this.generateNextBiography(i + 1);
      },
      generateBiography() {
        this.stage = STAGE.BIOGRAPHY_CREATION;
        this.generateNextBiography(0);
      },
      acceptBiography() {
        this.useBiography = true;
        this.finish();
      },
      finish() {
        this.stage = STAGE.APPLYING;
        const options = ExternalService.createRequestOptions(
          this.name,
          this.epitaphy,
          this.author,
          this.splitName(this.name).firstName,
          this.splitName(this.name).lastName,
          null,
          this.places[0],
          this.places[1],
          this.children,
          this.partner,
          this.citizenship,
          this.education,
          this.career,
          this.achievments,
          this.acceptBiography ? this.biography[0].header : "",
          this.acceptBiography ? this.biography[0].text : "",
          this.acceptBiography ? this.biography[1].header : "",
          this.acceptBiography ? this.biography[1].text : "",
          this.acceptBiography ? this.biography[2].header : "",
          this.acceptBiography ? this.biography[2].text : "",
          this.acceptBiography ? this.biographyEnding : "",
          this.dates[0],
          this.dates[1]
        );
        ExternalService.sendPage(options, (data) => {
          debugger;
        }, err => {
          this.$notify({text:"Не удалось загрузить страницу на сервер, попробуйте позже.", type: "error"});
          console.log(err);
        })
      },
      // UTILS
      splitName(name) {
        const nameParts = name.split(' ');
        if (nameParts.length === 1) return { firstName: nameParts[0], patronymic: null, lastName: null };
        if (nameParts.length === 2) return { firstName: nameParts[0], patronymic: null, lastName: nameParts[1] };
        if (nameParts.length === 3) return { firstName: nameParts[0], patronymic: nameParts[1], lastName: nameParts[2] };
        return { firstName: null, patronymic: null, lastName: null };
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