import { NetworkService } from "./NetworkService";

export const MemoryService = {
    getQuestions: (model, name, sex, birth_date, success, fail) => {
        NetworkService.ClassicRequest("POST", `v1/completion/${model}/questions`,
            {human_info: {name, sex, birth_date}}, response => {
            success(response.data);
        }, fail);
    },
    getBiographyQuestions: (model, type_of_story, name, sex, birth_date, success, fail) => {
        NetworkService.ClassicRequest("POST", `v1/completion/${model}/questions/biography`,
            {type_of_story, human_info: {name, sex, birth_date}}, response => {
            success(response.data);
        }, fail);
    },
    getEpitaph: (model, name, sex, birth_date, death_date, questions, success, fail) => {
        NetworkService.ClassicRequest("POST", `v1/completion/${model}/epitaph`,
            {human_info: {name, sex, birth_date, death_date, questions}}, response => {
            success(response.data);
        }, fail);
    },
    getBiography: (model, type_of_story, name, sex, birth_date, death_date, questions, previous, success, fail) => {
        NetworkService.ClassicRequest("POST", `v1/completion/${model}/biography`,
            {type_of_story, human_info: {name, sex, birth_date, death_date, questions}, previous}, response => {
            success(response.data);
        }, fail);
    },
    getBiographyShort: (model, biography, name, sex, birth_date, death_date, success, fail) => {
        NetworkService.ClassicRequest("POST", `v1/completion/${model}/biography/short`,
            {part_young: biography[0].text, part_middle: biography[1].text, part_old: biography[2].text, human_info: {name, sex, birth_date, death_date}}, response => {
            success(response.data);
        }, fail);
    },
}