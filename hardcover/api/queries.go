package api

const authorQuery = `
query GetFilterListWithAuthorBooks($id: Int!, $limit: Int!, $offset: Int!, $minUsersCount: Int!) {
    filter_lists(where: { id: { _eq: $id } }) {
        authors_count
        id
        updated_at
        filter_list_entities {
            author {
                alternate_names
                id
                name
                object_type
                slug
                state
                title
                user_id
                users_count
                books_count
                contributions(
                    where: {
                        contributable_type: { _eq: "Book" }
                        book: {
                            users_read_count: { _gt: $minUsersCount }
                            editions: { language_id: { _eq: 1 } }
                        }
                    }
                    limit: $limit
                    offset: $offset
                ) {
                    book {
                        id
                        title
                        literary_type_id
                        book_category_id
                        slug
                        subtitle
                        featured_book_series_id
                        activities_count
                        users_count
                        book_series {
                            id
                            position
                            series {
                                id
                                name
                            }
                        }
                        editions {
                            isbn_10
                            isbn_13
                            id
                        }
                    }
                }
            }
        }
    }
}

`

const toReadQuery = `
query WantToReadBooks($limit: Int!, $offset: Int!){
   me {
        user_books(
            where: {status_id: {_eq: 1}}
            limit: $limit
            offset: $offset
        ) {
            book {
                id
                title
                contributions {
                    author {
                        name
                    }
                } 
                literary_type_id
                book_category_id
                slug
                subtitle
                featured_book_series_id
                activities_count
                users_count
                book_series {
                    id
                    position
                    series {
                        id
                        name
                    }
                }
                editions {
                    isbn_10
                    isbn_13
                    id
                }
            }
        }
    }
}

`
