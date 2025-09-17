const User = require('./models/user');

const register = async (user) => {

    try { 
        const { name, email, password } = user;
        const query = {name: {$eq: name}, email: {$eq: email}, password: {$eq: password}}
        
        const existUser = await User.findOne({email: query.email});

        if(existUser) { return null }

        const newUser = new User({
            name: name,
            email: email,
            password: password
        });

        await newUser.save();

        return newUser;
    }

    catch(error) { throw error; }
    
}

const login = async (credentials) => {

    try {
        const { email, password } = credentials;
        const query = {email: {$eq: credentials.email}, password: {$eq: credentials.password}}

        const existsUser = await User.find({$and: [ { email: query.email}, { password: query.password} ]});

        if(!existsUser) { return null;}

        const returnUser = existsUser.map((user) => {
            return user.email
        })


        return returnUser;
    }

    catch(error) { throw error; }
    

}

module.exports = {
    register,
    login
};