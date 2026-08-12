package api

const authorQuery = `
{
    filter_lists(where: { id: { _eq: --TOFOLLOWED_AUTHORS_LIST_ID-- } }) {
        authors_count
        filter_list_entities {
            created_at
            entity_id
            entity_type
            filter_list_id
            id
            updated_at
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
                contributions(where: { contributable_type: { _eq: "Book" } }) {
                    book {
                        id
                        title
                        activities_count
                        users_count
                        book_series {
                            series {
                                name
                            }
                        }
                    }
                }
            }
           
        }
    }
}
`

const toReadQuery = `
query {
			me {
				user_books(where: {status_id: {_eq: 1}}) {
					book {
						id
						title
                        contributions {
							author {
								name
							}
						}
					}
				}
			}
		}
`
