package db

import (
	"context"
	"database/sql"

	"fmt"
	"log"
	"math/rand"
	"socialone/internal/store"
)

var usernames = []string{
    "alice", "bob", "dave", "eve", "frank", "grace", "heidi", "ivan", "judy", "karl",
    "leo", "mike", "nancy", "olivia", "paul", "quinn", "rachel", "sam", "tina", "uma",
    "vicky", "walter", "xander", "yara", "zane", "amber", "brad", "claire", "danny",
    "ella", "felix", "gina", "hank", "ian", "julia", "kevin", "lily", "matt", "nina",
    "oscar", "peter", "quincy", "ronnie", "susan", "tom", "ursula", "victor", "will",
    "xena", "yasmine", "zach", "andrea", "ben", "chris",
}

var titles = []string{
    "5 Tips for Better Productivity",
    "How to Stay Motivated Every Day",
    "The Future of Remote Work",
    "Beginner's Guide to Mindfulness",
    "Top 10 Travel Destinations in 2025",
    "Healthy Habits for a Better Life",
    "The Power of Positive Thinking",
    "How to Build a Personal Brand Online",
    "Why Reading is Good for Your Brain",
    "Simple Ways to Boost Your Creativity",
    "Mastering Time Management",
    "How to Make Your First Million",
    "Essential Tools for Entrepreneurs",
    "What to Do When You Feel Burned Out",
    "How to Start a Blog in 2025",
    "The Art of Saying No",
    "How to Grow Your Instagram Following",
    "Top Strategies for Self-Care",
    "How to Build a Morning Routine",
    "Secrets to Staying Fit Without the Gym",
}


var contents = []string{
    "Productivity isn't just about working harder, it's about working smarter. Here are 5 tips that can help you make the most of your time and focus on the tasks that matter most.",
    "Staying motivated every day can be tough, but it's possible with the right mindset. Here are a few strategies to help you maintain your energy and stay on track.",
    "Remote work is here to stay, but how will it evolve? Let's explore how remote work is changing the workplace and what that means for the future.",
    "Mindfulness is all about being present. This beginner's guide will show you how to practice mindfulness in your daily life and reduce stress along the way.",
    "From stunning beaches to vibrant cities, 2025 has plenty of exciting travel destinations. Discover the top 10 places to visit this year, each with its own unique appeal.",
    "Building healthy habits can transform your life. In this post, we’ll discuss simple changes that can help improve your physical, mental, and emotional well-being.",
    "Positive thinking isn't just a buzzword—it’s a mindset. Learn how focusing on the positives can lead to a happier and more productive life.",
    "Building a personal brand requires consistency and strategy. Here’s how you can start creating an online presence that reflects your values and expertise.",
    "Reading is an excellent habit that improves your focus, creativity, and cognitive function. Here's why reading is important and how to make it a daily habit.",
    "Creativity isn't something you’re born with; it’s something you can cultivate. Try these simple techniques to jump-start your creative process and boost your ideas.",
    "Time management is all about maximizing your productivity while keeping your sanity intact. In this post, I’ll share my favorite time management techniques to help you make the most of your day.",
    "Building wealth doesn’t happen overnight, but with the right mindset, it’s possible. Here are some essential steps to take toward your first million dollars.",
    "Entrepreneurs need the right tools to succeed. In this post, I’ll go over the essential software and apps you need to streamline your workflow and grow your business.",
    "Burnout is real and can affect anyone. Learn how to recognize the signs of burnout and take steps to recover and prevent it from happening in the future.",
    "Starting a blog in 2025 is easier than ever. Learn the essential steps for setting up your blog, creating content, and growing your audience from scratch.",
    "Saying no is an essential skill that can help you protect your time and energy. Here’s how to do it gracefully while still maintaining strong relationships.",
    "Instagram is a powerful platform for building your brand. Discover proven strategies to grow your Instagram following and connect with your audience in a meaningful way.",
    "Self-care is essential for maintaining mental and physical health. In this post, I’ll share some simple yet effective self-care strategies to help you feel your best every day.",
    "Starting your day off on the right foot can set the tone for the rest of your day. Here’s how to design a morning routine that energizes you and sets you up for success.",
    "You don’t need a gym to stay fit. Here are some fun and effective ways to stay in shape without ever stepping foot inside a gym.",
}


var tags = []string{
    "productivity", "mindfulness", "motivation", "self-care", "entrepreneurship",
    "travel", "creativity", "wellness", "time-management", "personal-growth",
    "business", "health", "inspiration", "digital-marketing", "startup",
    "fitness", "positive-thinking", "blogging", "mental-health", "self-improvement",
}




var comments= []string{
    "This is such a helpful post, thanks for sharing these tips!",
    "I never thought about that perspective before. Very insightful!",
    "Great read, I will definitely try these suggestions.",
    "I love how simple yet effective these strategies are. Thanks!",
    "This post came at the perfect time, I really needed this advice.",
    "I’m going to implement these ideas in my routine right away!",
    "I’ve been struggling with motivation lately, this really helped!",
    "Such a great guide for beginners, I feel more confident now.",
    "The way you explain things is so clear and easy to follow, thank you.",
    "I love how practical these tips are. I’ll be bookmarking this post!",
    "I’m always looking for new productivity hacks, this one is gold!",
    "This is exactly what I needed to read today, thank you for sharing.",
    "These tips are so simple but powerful. I’m excited to try them out!",
    "Your advice really resonates with me, I’ve been feeling overwhelmed.",
    "I shared this post with my friend, they will love it too!",
    "This blog has been a game-changer for me, keep the posts coming!",
    "I’m so glad I found this post! It’s exactly what I was looking for.",
    "The content is amazing, I’ll definitely be following your blog for more tips.",
    "This was such an eye-opener! I can’t wait to put these ideas into action.",
    "Your posts always bring a fresh perspective, I appreciate your insights.",
}


func Seed(store store.Storage, db *sql.DB){
	ctx:= context.Background()

	users:= generateUsers(100)
    tx,_:= db.BeginTx(ctx,nil)
	for _,user := range users{
		if err:= store.Users.Create(ctx,tx,user); err!=nil{
            _=tx.Rollback()
			log.Println("Error creating user:", err)
            return

		}
	}
    tx.Commit()

	posts:= generatePost(200,users)
	for _,post:=range posts{
		if err:= store.Posts.Create(ctx,post);err!=nil{
			log.Println("Error creating the post",err)
			return
		}

	}

	comments:= generateComment(500,users,posts)
	for _,comment:=range comments {
		if err:= store.Comments.Create(ctx,comment);err!=nil{
			log.Println("Error creating the post",err)
			return
		}

	}

	log.Println("seeding complete")


}
func generateUsers(num int) [] *store.User{
	users:= make([]*store.User, num)



	for i:= 0;i<num; i++{
		users[i]= &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email: usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
            RoleID: store.Role{
                Name: "user",
            },
        }

		}



	}
	return users
}
func generatePost(nums int, users []*store.User) []*store.Post{
	posts:= make ([]*store.Post, nums)
	for i:=0;i<nums; i++ {
		user:= users[rand.Intn(len(users))]
		posts[i]= &store.Post{
			UserId: user.ID,
			Title: titles[rand.Intn(len(titles))],
			Content: titles[rand.Intn(len(contents))],
			Tags:  []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],

			},
		}
	}
	return posts
}
func generateComment(num int, users []*store.User, posts []*store.Post) []*store.Comments {
    cms := make([]*store.Comments, num)  // Initialize the slice to hold 'num' comments

    for i := 0; i < num; i++ {
        cms[i] = &store.Comments{
            PostID: posts[rand.Intn(len(posts))].ID,  // Random post
            UserID: users[rand.Intn(len(users))].ID,  // Random user
            Content: comments[rand.Intn(len(comments))], // Placeholder for dynamic comment content
        }
    }

    return cms
}