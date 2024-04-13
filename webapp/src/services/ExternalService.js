import { NetworkService } from "./NetworkService";

export const ExternalService = {
    createRequestOptions(name, epitaph,
        author_epitaph, firstName,
        lastName, surname,
        placeOfBirth, placeOfDeath,
        children, spouse,
        citizenship, education,
        occupation, awards,
        bio_title1, bio1_text,
        bio_title2, bio2_text,
        bio_title3, bio3_text, birth, death
    ) {
        function formatDateToString(date) {
            const year = date.getFullYear();
            const month = ('0' + (date.getMonth() + 1)).slice(-2);
            const day = ('0' + date.getDate()).slice(-2);
            const hours = ('0' + date.getHours()).slice(-2);
            const minutes = ('0' + date.getMinutes()).slice(-2);
            const seconds = ('0' + date.getSeconds()).slice(-2);
            const formattedDate = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
            return formattedDate;
        }
        return {
            method: 'PUT',
            url: 'https://mc.dev.rand.agency/api/page/54553748',
            headers: {
                'User-Agent': 'Thunder Client (https://www.thunderclient.com)',
                Accept: 'application/json',
                'Content-Type': 'application/json;charset=UTF-8',
                Authorization: 'Bearer 4221|awGMegjL4NXWsFB9QfunXGvqDSkLwlDHN6KbyJK0'
            },
            data: {
                id: 8736,
                name: name,
                surname: surname,
                patronym: null,
                birthday_at: (formatDateToString(birth)),
                died_at: (formatDateToString(death)),
                epitaph: epitaph,
                author_epitaph: author_epitaph,
                page_type_id: 1,
                slug: 54553748,
                filled_fields: ['biography_1', 'biography_2', 'end_of_biography', 'epitaph'],
                lastName: lastName,
                firstName: firstName,
                link: 'https://mc.dev.rand.agency/page/54553748',
                full_name: lastName + ' ' + firstName + ' ' + surname,
                page_type_name: 'pageType.full',
                biographies: [
                    {
                        id: null,
                        title: bio_title1,
                        description: bio1_text,
                        page_id: 8736,
                        checked: true
                    },
                    {
                        id: null,
                        title: bio_title2,
                        description: bio2_text,
                        page_id: 8736,
                        checked: true
                    },
                    {
                        id: null,
                        title: bio_title3,
                        description: bio3_text,
                        page_id: 8736,
                        checked: true
                    }
                ],
                page_information: [
                    { page_id: 8736, title: 'pageInformation.placeOfBirth', description: placeOfBirth },
                    { page_id: 8736, title: 'pageInformation.placeOfDeath', description: placeOfDeath },
                    { page_id: 8736, title: 'pageInformation.children', description: children },
                    { page_id: 8736, title: 'pageInformation.wife||pageInformation.husband', description: spouse },
                    { page_id: 8736, title: 'pageInformation.citizenship', description: citizenship },
                    { page_id: 8736, title: 'pageInformation.education', description: education },
                    { page_id: 8736, title: 'pageInformation.occupation', description: occupation },
                    { page_id: 8736, title: 'pageInformation.awards', description: awards }
                ]
            }
        };
    },
    sendPage(requestOptions, success, fail) {
        return NetworkService.RawRequest(requestOptions, response => {
            success(response.data);
        }, fail);
    }
}