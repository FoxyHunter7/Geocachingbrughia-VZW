<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Message extends Model
{
    use HasFactory;

    protected $hidden = [
        'id',
        'created_at'
    ];
    protected $fillable = [
        'state'
    ];

    public function translations()
    {
        return $this->hasMany(MessageLanguage::class, 'message_id', 'id');
    }
}
