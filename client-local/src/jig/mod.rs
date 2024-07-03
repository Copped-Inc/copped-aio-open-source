use std::collections::HashMap;
use std::ops::Add;
use rand::Rng;

pub(crate) fn address(mut s: String) -> String {

    s = xyz().to_string().add(" ").add(s.as_str()).add(" ").add(xyz().as_str()).to_string();
    s = s.replace("  ", " ");

    let mut rng = rand::thread_rng();
    for _ in 0..2 {
        let pos = rng.gen_range(0..s.len());
        s.insert(pos, random_char());
    }

    s

}

pub(crate) fn name(mut s: String) -> String {

    let mut rng = rand::thread_rng();
    for _ in 0..2 {
        let pos = rng.gen_range(0..s.len());
        s.insert(pos, random_char());
    }

    s

}

pub(crate) fn xyz() -> String {

    let mut s = String::new();
    for _ in 0..3 {
        s.push(random_char());
    }

    s

}

pub(crate) fn random_char() -> char {

    let chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";
    let mut rng = rand::thread_rng();
    chars.chars().nth(rng.gen_range(0..chars.len())).unwrap()

}

pub(crate) fn random_number(length: i32) -> String {

    let mut rng = rand::thread_rng();
    let mut s = String::new();
    for _ in 0..length {
        s.push(rng.gen_range(0..10).to_string().parse().unwrap());
    }

    s

}

pub(crate) fn random_name() -> String {

    let names: [&str; 200] = [
        "Michael",
        "Christopher",
        "Jessica",
        "Matthew",
        "Ashley",
        "Jennifer",
        "Joshua",
        "Amanda",
        "Daniel",
        "David",
        "James",
        "Robert",
        "John",
        "Joseph",
        "Andrew",
        "Ryan",
        "Brandon",
        "Jason",
        "Justin",
        "Sarah",
        "William",
        "Jonathan",
        "Stephanie",
        "Brian",
        "Nicole",
        "Nicholas",
        "Anthony",
        "Heather",
        "Eric",
        "Elizabeth",
        "Adam",
        "Megan",
        "Melissa",
        "Kevin",
        "Steven",
        "Thomas",
        "Timothy",
        "Christina",
        "Kyle",
        "Rachel",
        "Laura",
        "Lauren",
        "Amber",
        "Brittany",
        "Danielle",
        "Richard",
        "Kimberly",
        "Jeffrey",
        "Amy",
        "Crystal",
        "Michelle",
        "Tiffany",
        "Jeremy",
        "Benjamin",
        "Mark",
        "Emily",
        "Aaron",
        "Charles",
        "Rebecca",
        "Jacob",
        "Stephen",
        "Patrick",
        "Sean",
        "Erin",
        "Zachary",
        "Jamie",
        "Kelly",
        "Samantha",
        "Nathan",
        "Sara",
        "Dustin",
        "Paul",
        "Angela",
        "Tyler",
        "Scott",
        "Katherine",
        "Andrea",
        "Gregory",
        "Erica",
        "Mary",
        "Travis",
        "Lisa",
        "Kenneth",
        "Bryan",
        "Lindsey",
        "Kristen",
        "Jose",
        "Alexander",
        "Jesse",
        "Katie",
        "Lindsay",
        "Shannon",
        "Vanessa",
        "Courtney",
        "Christine",
        "Alicia",
        "Cody",
        "Allison",
        "Bradley",
        "Samuel",
        "Shawn",
        "April",
        "Derek",
        "Kathryn",
        "Kristin",
        "Chad",
        "Jenna",
        "Tara",
        "Maria",
        "Krystal",
        "Jared",
        "Anna",
        "Edward",
        "Julie",
        "Peter",
        "Holly",
        "Marcus",
        "Kristina",
        "Natalie",
        "Jordan",
        "Victoria",
        "Jacqueline",
        "Corey",
        "Keith",
        "Monica",
        "Juan",
        "Donald",
        "Cassandra",
        "Meghan",
        "Joel",
        "Shane",
        "Phillip",
        "Patricia",
        "Brett",
        "Ronald",
        "Catherine",
        "George",
        "Antonio",
        "Cynthia",
        "Stacy",
        "Kathleen",
        "Raymond",
        "Carlos",
        "Brandi",
        "Douglas",
        "Nathaniel",
        "Ian",
        "Craig",
        "Brandy",
        "Alex",
        "Valerie",
        "Veronica",
        "Cory",
        "Whitney",
        "Gary",
        "Derrick",
        "Philip",
        "Luis",
        "Diana",
        "Chelsea",
        "Leslie",
        "Caitlin",
        "Leah",
        "Natasha",
        "Erika",
        "Casey",
        "Latoya",
        "Erik",
        "Dana",
        "Victor",
        "Brent",
        "Dominique",
        "Frank",
        "Brittney",
        "Evan",
        "Gabriel",
        "Julia",
        "Candice",
        "Karen",
        "Melanie",
        "Adrian",
        "Stacey",
        "Margaret",
        "Sheena",
        "Wesley",
        "Vincent",
        "Alexandra",
        "Katrina",
        "Bethany",
        "Nichole",
        "Larry",
        "Jeffery",
        "Curtis",
        "Carrie",
        "Todd",
        "Blake",
        "Christian",
        "Randy",
        "Dennis",
        "Alison",
    ];

    let mut rng = rand::thread_rng();
    let i = rng.gen_range(0..names.len());
    let name = names[i];
    name.to_string()

}

pub(crate) fn country_id(s: String) -> String {

    let mut m: HashMap<&str, &str> = HashMap::new();
    m.insert("AL", "7");
    m.insert("AD", "2");
    m.insert("AT", "13");
    m.insert("BE", "21");
    m.insert("BA", "18");
    m.insert("BG", "23");
    m.insert("CA", "37");
    m.insert("HR", "97");
    m.insert("CY", "53");
    m.insert("CZ", "54");
    m.insert("DK", "57");
    m.insert("EE", "62");
    m.insert("FI", "247");
    m.insert("FR", "74");
    m.insert("DE", "69");
    m.insert("GI", "82");
    m.insert("GR", "88");
    m.insert("GG", "80");
    m.insert("HU", "99");
    m.insert("IS", "108");
    m.insert("IE", "101");
    m.insert("IM", "103");
    m.insert("IT", "91");
    m.insert("JE", "110");
    m.insert("LV", "134");
    m.insert("LI", "128");
    m.insert("LT", "132");
    m.insert("LU", "133");
    m.insert("MT", "151");
    m.insert("MD", "138");
    m.insert("MC", "137");
    m.insert("ME", "139");
    m.insert("NL", "164");
    m.insert("MK", "142");
    m.insert("NO", "165");
    m.insert("PL", "177");
    m.insert("PT", "182");
    m.insert("RO", "157");
    m.insert("SM", "202");
    m.insert("RS", "188");
    m.insert("SK", "200");
    m.insert("SI", "198");
    m.insert("ES", "66");
    m.insert("SE", "195");
    m.insert("CH", "190");
    m.insert("TR", "221");
    m.insert("GB", "211");
    m.insert("US", "230");
    m.insert("VA", "233");

    m.get(&s.as_str()).unwrap().to_string()

}