import { NetworkService } from "./NetworkService";

export const MemoryService = {
    getQuestions: (model, name, sex, birth_date, success, fail) => {
        NetworkService.ClassicRequest("POST", `v1/completion/${model}/questions`,
            {human_info: {name, sex, birth_date}}, response => {
            success(response.data);
        }, fail);
    },
    getStory: (model, type_of_story, name, sex, birth_date, death_date, questions, success, fail) => {
        NetworkService.ClassicRequest("POST", `v1/completion/${model}/story`,
            {type_of_story, human_info: {name, sex, birth_date, death_date, questions}}, response => {
            success(response.data);
        }, fail);
    }
}