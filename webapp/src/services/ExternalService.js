import { NetworkService } from "./NetworkService";

export const ExternalService = {
    createRequestOptions(page, isMale, name, epitaph,
        author_epitaph, firstName,
        lastName, surname,
        placeOfBirth, placeOfDeath,
        children, spouse,
        citizenship, education,
        occupation, awards,
        bio_title1, bio1_text,
        bio_title2, bio2_text,
        bio_title3, bio3_text,
        bioend_text,
        birth, death
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
        function formatDate(date) {
            let day = String(date.getDate()).padStart(2, '0');
            let month = String(date.getMonth() + 1).padStart(2, '0');
            let year = date.getFullYear();
            return { day: day, month: month, year: year };
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
                id: page.id,
                name: name,
                surname: surname,
                patronym: null,
                birthday_at: (formatDateToString(birth)),
                died_at: (formatDateToString(death)),
                start: formatDate(birth),
                end: formatDate(death),
                epitaph: epitaph,
                author_epitaph: author_epitaph,
                page_type_id: 1,
                slug: page.slug,
                filled_fields: ['biography_1', 'biography_2', 'biography_3', 'end_of_biography', 'epitaph'],
                lastName: lastName,
                firstName: firstName,
                link: page.link,
                full_name: lastName + ' ' + firstName + ' ' + surname,
                page_type_name: 'pageType.full',
                biographies: [
                    {
                        id: null,
                        title: bio_title1,
                        description: bio1_text,
                        page_id: page.id,
                        checked: true
                    },
                    {
                        id: null,
                        title: bio_title2,
                        description: bio2_text,
                        page_id: page.id,
                        checked: true
                    },
                    {
                        id: null,
                        title: bio_title3,
                        description: bio3_text,
                        page_id: page.id,
                        checked: true
                    },
                    {
                        id: null,
                        title: "Заключение",
                        description: bioend_text,
                        page_id: page.id,
                        checked: true
                    }
                ],
                page_information: [
                    { page_id: page.id, title: 'pageInformation.placeOfBirth', description: placeOfBirth },
                    { page_id: page.id, title: 'pageInformation.placeOfDeath', description: placeOfDeath },
                    { page_id: page.id, title: 'pageInformation.children', description: children },
                    { page_id: page.id, title: isMale ? 'pageInformation.wife' : 'pageInformation.husband', description: spouse },
                    { page_id: page.id, title: 'pageInformation.citizenship', description: citizenship },
                    { page_id: page.id, title: 'pageInformation.education', description: education },
                    { page_id: page.id, title: 'pageInformation.occupation', description: occupation },
                    { page_id: page.id, title: 'pageInformation.awards', description: awards }
                ]
            }
        };
    },
    sendPage(requestOptions, success, fail) {
        return NetworkService.RawRequest(requestOptions, response => {
            success(response.data);
        }, fail);
    },
    getPages(success, fail, cookies) {
        return NetworkService.AuthRequest("GET", `cabinet/individual-pages`, null, cookies, response => {
            success(response.data);
        }, fail);
    },
    sendReport(page_id, fio, email, text, relation_role, success, fail) {
        return NetworkService.DirectRequest("POST", `comment`,
        { page_id, fio, email, text, relation_role, checked: true, hasEmail: email.length > 0},
        response => {
            success(response.data);
        }, fail);
    }
}